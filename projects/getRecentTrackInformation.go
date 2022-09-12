package projects

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

type GetRecentTrackInformation interface {
	Execute()
}

type getRecentTrackInformation struct {
	fromDate          int64
	lastFMHandler     endpoints.LastFMHandler
	tracksForDuration []models.Track
	spotifyHandler    endpoints.SpotifyHandler
	audioFeatures     map[spotify.ID]*spotify.AudioFeatures
}

func NewGetRecentTrackInformation(fromDate int64, lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) GetRecentTrackInformation {
	return &getRecentTrackInformation{
		fromDate:          fromDate,
		lastFMHandler:     lastFMHandler,
		tracksForDuration: make([]models.Track, 0),
		spotifyHandler:    spotifyHandler,
		audioFeatures:     models.GetSpotifyIDToAudioFeatures(),
	}
}

func (c *getRecentTrackInformation) Execute() {
	c.getInformation()
	c.printoutResultsToTxt()
	//c.printoutResultsWFHValenceComparison()
}

func (c *getRecentTrackInformation) getInformation() {
	c.tracksForDuration = c.lastFMHandler.GetAllRecentTracks(c.fromDate, c.lastFMHandler.GetUserInfo().Name)
	fmt.Println("-----------------------------")
	fmt.Println("There are this many tracks: ", len(c.tracksForDuration))
	nonCachedTrackIDs := make([]spotify.ID, 0)
	couldNotFindInSearch := 0
	couldNotMatchInSearch := 0
	trackToIDHash := models.GetSpotifySearchToSongIDs()
	for i := 0; i < len(c.tracksForDuration); {
		t := c.tracksForDuration[i]
		if i != 0 && i%5000 == 0 {
			fmt.Println("5 second search sleep")
			fmt.Println("Search index: ", i)
			time.Sleep(5 * time.Second)
		}

		searchKey := t.Artist + " " + t.AlbumName + " " + t.Name
		if value, ok := trackToIDHash[searchKey]; ok {
			if value != constants.NotFound {
				t.SpotifyID = value
				c.tracksForDuration[i] = t
			} else {
				t.SpotifyID = constants.NotFound
				c.tracksForDuration[i] = t
			}
			i += 1
			continue
		}

		searchResult, err := c.spotifyHandler.SearchForSong(t.Artist, t.AlbumName, t.Name)
		if err != nil {
			fmt.Println("5 second search error sleep")
			fmt.Println("Search error index: ", i)
			time.Sleep(5 * time.Second)
			continue
		}
		if searchResult != nil {
			comparisonResult := compareMultipleReturnedTracks(t, searchResult)
			if comparisonResult != constants.EmptyString {
				nonCachedTrackIDs = append(nonCachedTrackIDs, comparisonResult)
				t.SpotifyID = comparisonResult
				c.tracksForDuration[i] = t
				trackToIDHash[searchKey] = comparisonResult
			} else {
				t.SpotifyID = constants.NotFound
				trackToIDHash[searchKey] = constants.NotFound
				c.tracksForDuration[i] = t
				couldNotMatchInSearch += 1
			}
		} else {
			t.SpotifyID = constants.NotFound
			trackToIDHash[searchKey] = constants.NotFound
			c.tracksForDuration[i] = t
			couldNotFindInSearch += 1
		}
		i += 1
	}
	models.AddSpotifySearchToSongIDs(trackToIDHash)

	fmt.Println("-----------------------------")
	fmt.Println("Could not match in search: ", couldNotMatchInSearch)
	fmt.Println("Could not find in search: ", couldNotFindInSearch)
	fmt.Println("Total tracks: ", len(c.tracksForDuration))
	fmt.Println("10 refresh sleep")
	time.Sleep(10 * time.Second)
	fmt.Println("-----------------------------")

	for i := 0; i < len(nonCachedTrackIDs); {
		upperLimit := i + 50
		if upperLimit > len(nonCachedTrackIDs) {
			upperLimit = len(nonCachedTrackIDs)
		}
		audioFeatures := c.spotifyHandler.GetAudioFeaturesOfTrack(nonCachedTrackIDs[i:upperLimit])
		for _, a := range audioFeatures {
			if _, ok := c.audioFeatures[a.ID]; !ok {
				c.audioFeatures[a.ID] = a
			}
		}
		if i != 0 && i%1000 == 0 {
			fmt.Println("5 second audio feature sleep")
			fmt.Println("AudioFeatures index: ", i)
			time.Sleep(5 * time.Second)
		}
		i += 50
	}
	if _, ok := c.audioFeatures[constants.NotFound]; !ok {
		c.audioFeatures[constants.NotFound] = nil
	}
	models.AddSpotifyIDToAudioFeatures(c.audioFeatures)

	fmt.Println("-----------------------------")
	fmt.Println("Searched for this many track IDs: ", len(nonCachedTrackIDs))
}

func (c *getRecentTrackInformation) printoutResultsToTxt() {
	f, err := os.Create("audioFeaturesTimeSeries_" + constants.Now.Format("20060102150405") + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	dataWriter := bufio.NewWriter(f)

	_, _ = dataWriter.WriteString("ListenDate\tTrack\tAlbum\tArtist\tDuration(S)\tSpotifyID\tAcousticness\tDanceability\tEnergy\tInstrumentalness\tLiveness\tLoudness\tSpeechiness\tTempo\tValence\n")
	for _, t := range c.tracksForDuration {
		af := c.audioFeatures[t.SpotifyID]

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

// Result
func (c *getRecentTrackInformation) printoutResultsWFHValenceComparison() {

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

	for _, t := range c.tracksForDuration {
		// I only want dem weekday tracks
		if t.ListenDate.Weekday() == 0 || t.ListenDate.Weekday() == 6 {
			continue
		}

		// I only want dem work hour tracks 10 AM - 6 PM
		if t.ListenDate.Hour() < 10 || t.ListenDate.Hour() > 18 {
			continue
		}

		// An audio feature id is required
		if af, ok := c.audioFeatures[t.SpotifyID]; !ok || af == nil {
			continue
		}

		wantedValue := int(c.audioFeatures[t.SpotifyID].Valence * 100.00)

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
