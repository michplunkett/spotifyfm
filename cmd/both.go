package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

func NewBothCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "both",
		Short: "Runs subcommands based on LastFM and Spotify data.",
	}
}

func NewRecentTrackInformationCmd(lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recent-track-info",
		Short: "Gets recent track information by combining information from LastFM and Spotify.",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting the recent track fetching process...")
			tracksForDuration, audioFeatures := getRecentTrackInformation(constants.StartOf2023, lastFMHandler, spotifyHandler)
			printoutResultsToTxt(tracksForDuration, audioFeatures)
		},
	}

	return cmd
}

func getRecentTrackInformation(fromDate int64, lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) ([]models.Track, map[spotify.ID]*spotify.AudioFeatures) {
	tracksForDuration := lastFMHandler.GetAllRecentTracks(fromDate, lastFMHandler.GetUserInfo().Name)
	models.AddLastFMTrackList(tracksForDuration)
	fmt.Println("-----------------------------")
	fmt.Printf("There are this many tracks: %d\n", len(tracksForDuration))
	nonCachedTrackIDs := make([]spotify.ID, 0)
	couldNotFindInSearch := 0
	couldNotMatchInSearch := 0

	trackToIDHash := models.GetSpotifySearchToSongIDs()
	for i := 0; i < len(tracksForDuration); {
		t := tracksForDuration[i]
		if i != 0 && i%1000 == 0 {
			sleepPrint(5, "Spotify search to ID")
			fmt.Printf("Search index: %d\n", i)
		}

		searchKey := t.Artist + " " + t.AlbumName + " " + t.Name
		if value, ok := trackToIDHash[searchKey]; ok {
			if value != constants.NotFound {
				t.SpotifyID = value
				tracksForDuration[i] = t
			} else {
				t.SpotifyID = constants.NotFound
				tracksForDuration[i] = t
			}
			i += 1
			continue
		}

		searchResult, err := spotifyHandler.SearchForSong(t.Artist, t.AlbumName, t.Name)
		if err != nil {
			sleepPrint(10, "Spotify song search")
			fmt.Printf("Search error index: %d\n", i)
			continue
		}
		if searchResult != nil {
			comparisonResult := compareMultipleReturnedTracks(t, searchResult)
			if comparisonResult != constants.EmptyString {
				nonCachedTrackIDs = append(nonCachedTrackIDs, comparisonResult)
				t.SpotifyID = comparisonResult
				tracksForDuration[i] = t
				trackToIDHash[searchKey] = comparisonResult
			} else {
				t.SpotifyID = constants.NotFound
				trackToIDHash[searchKey] = constants.NotFound
				tracksForDuration[i] = t
				couldNotMatchInSearch += 1
			}
		} else {
			t.SpotifyID = constants.NotFound
			trackToIDHash[searchKey] = constants.NotFound
			tracksForDuration[i] = t
			couldNotFindInSearch += 1
		}
		i += 1
	}
	models.AddSpotifySearchToSongIDs(trackToIDHash)

	fmt.Println("-----------------------------")
	fmt.Printf("Could not match in search: %d\n", couldNotMatchInSearch)
	fmt.Printf("Could not find in search: %d\n", couldNotFindInSearch)
	fmt.Printf("Total tracks: %d\n", len(tracksForDuration))
	sleepPrint(10, "break between hitting the APIs")
	fmt.Println("-----------------------------")

	audioFeatures := make(map[spotify.ID]*spotify.AudioFeatures)
	for i := 0; i < len(nonCachedTrackIDs); {
		var upperLimit int
		if len(nonCachedTrackIDs) < i+50 {
			upperLimit = len(nonCachedTrackIDs)
		} else {
			upperLimit = i + 50
		}

		if upperLimit > len(nonCachedTrackIDs) {
			upperLimit = len(nonCachedTrackIDs)
		}
		features := spotifyHandler.GetAudioFeaturesOfTrack(nonCachedTrackIDs[i:upperLimit])
		for _, a := range features {
			if _, ok := audioFeatures[a.ID]; !ok {
				audioFeatures[a.ID] = a
			}
		}
		if i != 0 && i%1000 == 0 {
			sleepPrint(30, "Spotify audio feature search")
			fmt.Println("AudioFeatures index: ", i)
		}
		i += 50
	}
	if _, ok := audioFeatures[constants.NotFound]; !ok {
		audioFeatures[constants.NotFound] = nil
	}
	models.AddSpotifyIDToAudioFeatures(audioFeatures)

	fmt.Println("-----------------------------")
	fmt.Println("Searched for this many track IDs: ", len(nonCachedTrackIDs))

	return tracksForDuration, audioFeatures
}

func printoutResultsToTxt(tracks []models.Track, features map[spotify.ID]*spotify.AudioFeatures) {
	f, err := os.Create("audioFeaturesTimeSeries_" + constants.Now.Format("20060102150405") + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	dataWriter := bufio.NewWriter(f)

	_, _ = dataWriter.WriteString("ListenDate\tTrack\tAlbum\tArtist\tDuration(S)\tSpotifyID\tAcousticness\tDanceability\tEnergy\tInstrumentalness\tLiveness\tLoudness\tSpeechiness\tTempo\tValence\n")
	for _, t := range tracks {
		af := features[t.SpotifyID]

		if af == nil {
			continue
		}

		trackStringArray := make([]string, 0)
		// Listen Date
		listenDate := t.ListenDate.Format(time.RFC3339)
		trackStringArray = append(trackStringArray, listenDate)
		// Track
		trackStringArray = append(trackStringArray, t.Name)
		// Album
		trackStringArray = append(trackStringArray, t.AlbumName)
		// Artist
		trackStringArray = append(trackStringArray, t.Artist)
		// Duration
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", float64(af.Duration)/1000.00))
		// SpotifyID
		trackStringArray = append(trackStringArray, string(af.ID))
		// Acousticness
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Acousticness))
		// Danceability
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Danceability))
		// Energy
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Energy))
		// Instrumentalness
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Instrumentalness))
		// Liveness
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Liveness))
		// Loudness
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Loudness))
		// Speechiness
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Speechiness))
		// Tempo
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Tempo))
		// Valence
		trackStringArray = append(trackStringArray, fmt.Sprintf("%f", af.Valence))

		_, _ = dataWriter.WriteString(strings.Join(trackStringArray[:], "\t") + "\n")
	}

	_ = dataWriter.Flush()
}

func printoutResultsWFHValenceComparison(tracks []models.Track, features map[spotify.ID]*spotify.AudioFeatures) {

	// Calculation vars -- pre pandemic
	prePandemicMondayVarSum := 0
	prePandemicMondayTracks := 0
	prePandemicTuesdayVarSum := 0
	prePandemicTuesdayTracks := 0
	prePandemicWednesdayVarSum := 0
	prePandemicWednesdayTracks := 0
	prePandemicThursdayVarSum := 0
	prePandemicThursdayTracks := 0
	prePandemicFridayVarSum := 0
	prePandemicFridayTracks := 0

	// Calculation vars -- pandemic
	pandemicMondayVarSum := 0
	pandemicMondayTracks := 0
	pandemicTuesdayVarSum := 0
	pandemicTuesdayTracks := 0
	pandemicWednesdayVarSum := 0
	pandemicWednesdayTracks := 0
	pandemicThursdayVarSum := 0
	pandemicThursdayTracks := 0
	pandemicFridayVarSum := 0
	pandemicFridayTracks := 0

	for _, t := range tracks {
		// I only want dem weekday tracks
		if t.ListenDate.Weekday() == 0 || t.ListenDate.Weekday() == 6 {
			continue
		}

		// I only want dem work hour tracks 10 AM - 6 PM
		if t.ListenDate.Hour() < 10 || t.ListenDate.Hour() > 18 {
			continue
		}

		// An audio feature id is required
		if af, ok := features[t.SpotifyID]; !ok || af == nil {
			continue
		}

		wantedValue := int(features[t.SpotifyID].Valence * 100.00)

		if t.ListenDate.Before(constants.WFHStartDay) {
			if t.ListenDate.Weekday() == 1 {
				prePandemicMondayVarSum += wantedValue
				prePandemicMondayTracks++
			} else if t.ListenDate.Weekday() == 2 {
				prePandemicTuesdayVarSum += wantedValue
				prePandemicTuesdayTracks++
			} else if t.ListenDate.Weekday() == 3 {
				prePandemicWednesdayVarSum += wantedValue
				prePandemicWednesdayTracks++
			} else if t.ListenDate.Weekday() == 4 {
				prePandemicThursdayVarSum += wantedValue
				prePandemicThursdayTracks++
			} else if t.ListenDate.Weekday() == 5 {
				prePandemicFridayVarSum += wantedValue
				prePandemicFridayTracks++
			}
		} else {
			if t.ListenDate.Weekday() == 1 {
				pandemicMondayVarSum += wantedValue
				pandemicMondayTracks++
			} else if t.ListenDate.Weekday() == 2 {
				pandemicTuesdayVarSum += wantedValue
				pandemicTuesdayTracks++
			} else if t.ListenDate.Weekday() == 3 {
				pandemicWednesdayVarSum += wantedValue
				pandemicWednesdayTracks++
			} else if t.ListenDate.Weekday() == 4 {
				pandemicThursdayVarSum += wantedValue
				pandemicThursdayTracks++
			} else if t.ListenDate.Weekday() == 5 {
				pandemicFridayVarSum += wantedValue
				pandemicFridayTracks++
			}
		}
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithColorsOpts(opts.Colors{"#bbbbbb", "#88c442"}),
		charts.WithTitleOpts(
			opts.Title{Title: "Valence of Tracks During the Workday: Before and During the Pandemic"},
		),
	)

	bar.SetXAxis([]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}).
		AddSeries(
			"Pre-pandemic",
			[]opts.BarData{
				{
					Value: prePandemicMondayVarSum / prePandemicMondayTracks,
				},
				{
					Value: prePandemicTuesdayVarSum / prePandemicTuesdayTracks,
				},
				{
					Value: prePandemicWednesdayVarSum / prePandemicWednesdayTracks,
				},
				{
					Value: prePandemicThursdayVarSum / prePandemicThursdayTracks,
				},
				{
					Value: prePandemicFridayVarSum / prePandemicFridayTracks,
				},
			}).
		AddSeries(
			"Pandemic",
			[]opts.BarData{
				{
					Value: pandemicMondayVarSum / pandemicMondayTracks,
				},
				{
					Value: pandemicTuesdayVarSum / pandemicTuesdayTracks,
				},
				{
					Value: pandemicWednesdayVarSum / pandemicWednesdayTracks,
				},
				{
					Value: pandemicThursdayVarSum / pandemicThursdayTracks,
				},
				{
					Value: pandemicFridayVarSum / pandemicFridayTracks,
				},
			})

	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func compareMultipleReturnedTracks(localTrack models.Track, searchTracks []spotify.FullTrack) spotify.ID {
	for _, searchTrack := range searchTracks {
		if models.RemoveNonWordCharacters(searchTrack.Name) == models.RemoveNonWordCharacters(localTrack.Name) &&
			models.RemoveNonWordCharacters(searchTrack.Artists[0].Name) == localTrack.LowerCaseArtist {
			return searchTrack.ID
		}
	}
	return constants.EmptyString
}

func sleepPrint(duration int, message string) {
	fmt.Printf("%d second %s sleep\n", duration, message)
	time.Sleep(time.Duration(duration) * time.Second)
}
