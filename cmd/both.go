package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

const MaxRetries = 6

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
			tracksForDuration, audioFeatures := getRecentTrackInformation(constants.StartOf2024, lastFMHandler, spotifyHandler)
			printoutResultsToTxt(tracksForDuration, audioFeatures)
		},
	}

	return cmd
}

func getRecentTrackInformation(fromDate time.Time, lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) ([]models.Track, map[spotify.ID]*spotify.AudioFeatures) {
	fmt.Printf("Starting the recent track fetching process with %s as the start date...\n", fromDate.Format("01-02-2006"))
	tracksForDuration := lastFMHandler.GetAllRecentTracks(fromDate.Unix(), lastFMHandler.GetUserInfo().Name)
	models.AddLastFMTrackList(tracksForDuration)
	fmt.Println("-----------------------------")
	fmt.Printf("There are this many tracks: %d\n", len(tracksForDuration))
	couldNotFindInSearch := 0
	couldNotMatchInSearch := 0
	spotifyAPICallForTrack := 0

	trackToIDHash := models.GetSpotifySearchToSongIDs()
	for i := 0; i < len(tracksForDuration); {
		t := tracksForDuration[i]
		spotifyAPICallForTrack += 1
		if i != 0 && i%10000 == 0 {
			sleepPrint(10, "Spotify search to ID")
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
			spotifyAPICallForTrack = 0
			i += 1
			continue
		}

		searchResult, err := spotifyHandler.SearchForSong(t.Artist, t.AlbumName, t.Name)
		if err != nil && spotifyAPICallForTrack < MaxRetries {
			fmt.Println(searchKey)
			fmt.Println(err.Error())
			sleepPrint(5, "Spotify song search")
			fmt.Printf("Search error index: %d\n", i)
			continue
		}
		if searchResult != nil {
			comparisonResult := compareMultipleReturnedTracks(t, searchResult)
			if comparisonResult != constants.EmptyString {
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

		spotifyAPICallForTrack = 0
		i += 1
	}
	models.AddSpotifySearchToSongIDs(trackToIDHash)

	fmt.Println("-----------------------------")
	fmt.Printf("Could not match in search: %d\n", couldNotMatchInSearch)
	fmt.Printf("Could not find in search: %d\n", couldNotFindInSearch)
	fmt.Printf("Total tracks: %d\n", len(tracksForDuration))
	fmt.Println("-----------------------------")

	audioFeatures := models.GetSpotifyIDToAudioFeatures()
	audioFeaturesForSearch := make([]spotify.ID, 0)
	for _, v := range trackToIDHash {
		if _, ok := audioFeatures[v]; !ok {
			audioFeaturesForSearch = append(audioFeaturesForSearch, v)
		}
	}

	fmt.Printf("This many tracks need audio features: %d\n", len(audioFeaturesForSearch))

	for i := 0; i < len(audioFeaturesForSearch); {
		var upperLimit int
		if len(audioFeaturesForSearch) < i+50 {
			upperLimit = len(audioFeaturesForSearch)
		} else {
			upperLimit = i + 50
		}

		features := spotifyHandler.GetAudioFeaturesOfTrack(audioFeaturesForSearch[i:upperLimit])
		for _, a := range features {
			if a == nil {
				continue
			}
			if _, ok := audioFeatures[a.ID]; !ok {
				audioFeatures[a.ID] = a
			}
		}
		if i != 0 && i%10000 == 0 {
			sleepPrint(5, "Spotify audio feature search")
			fmt.Println("AudioFeatures index: ", i)
		}
		i += 50
	}
	delete(audioFeatures, constants.NotFound)
	models.AddSpotifyIDToAudioFeatures(audioFeatures)

	fmt.Println("-----------------------------")
	fmt.Println("Searched for this many track IDs: ", len(audioFeaturesForSearch))

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
