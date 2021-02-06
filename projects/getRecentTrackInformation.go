package projects

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
		audioFeatures:     make(map[spotify.ID]*spotify.AudioFeatures, 0),
	}
}

func (getInfo *getRecentTrackInformation) Execute() {
	getInfo.getInformation()
	getInfo.printoutResults()
}

func (getInfo *getRecentTrackInformation) getInformation() {
	getInfo.tracksForDuration = getInfo.lastFMHandler.GetAllRecentTracks(getInfo.fromDate, getInfo.lastFMHandler.GetUserInfo().Name)
	fmt.Println("-----------------------------")
	fmt.Println("There are this many tracks: ", len(getInfo.tracksForDuration))
	trackIDs := make([]spotify.ID, 0)
	couldNotFindInSearch := 0
	couldNotMatchInSearch := 0
	trackToIDHash := make(map[string]spotify.ID, 0)
	for i := 0; i < len(getInfo.tracksForDuration); {
		t := getInfo.tracksForDuration[i]
		if i != 0 && i%1000 == 0 {
			fmt.Println("5 second search sleep")
			fmt.Println("Search index: ", i)
			time.Sleep(5 * time.Second)
		}

		searchKey := t.Artist + " " + t.AlbumName + " " + t.Name
		if value, ok := trackToIDHash[searchKey]; ok {
			if value != constants.NotFound {
				t.SpotifyID = value
				getInfo.tracksForDuration[i] = t
			} else {
				t.SpotifyID = constants.NotFound
				getInfo.tracksForDuration[i] = t
			}
			i += 1
			continue
		}

		searchResult, err := getInfo.spotifyHandler.SearchForSong(t.Artist, t.AlbumName, t.Name)
		if err != nil {
			fmt.Println("10 second search error sleep")
			fmt.Println("Search error index: ", i)
			time.Sleep(10 * time.Second)
			continue
		}
		if searchResult != nil {
			comparisonResult := compareMultipleReturnedTracks(t, searchResult)
			if comparisonResult != constants.EmptyString {
				trackIDs = append(trackIDs, comparisonResult)
				t.SpotifyID = comparisonResult
				getInfo.tracksForDuration[i] = t
				trackToIDHash[searchKey] = comparisonResult
			} else {
				//fmt.Println("--- Count not match in search ---")
				//fmt.Println("index: ", i)
				//fmt.Println("search string: " + searchKey)
				//fmt.Println(searchResult)
				t.SpotifyID = constants.NotFound
				trackToIDHash[searchKey] = constants.NotFound
				getInfo.tracksForDuration[i] = t
				couldNotMatchInSearch += 1
			}
		} else {
			//fmt.Println("--- Count not find in search ---")
			//fmt.Println("index: ", i)
			//fmt.Println("search string: " + searchKey)
			t.SpotifyID = constants.NotFound
			trackToIDHash[searchKey] = constants.NotFound
			getInfo.tracksForDuration[i] = t
			couldNotFindInSearch += 1
		}
		i += 1
	}
	fmt.Println("-----------------------------")
	fmt.Println("Could not match in search: ", couldNotMatchInSearch)
	fmt.Println("Could not find in search: ", couldNotFindInSearch)
	fmt.Println("Total tracks: ", len(getInfo.tracksForDuration))
	fmt.Println("10 refresh sleep")
	time.Sleep(10 * time.Second)
	fmt.Println("-----------------------------")

	for i := 0; i < len(trackIDs); {
		upperLimit := i + 50
		if upperLimit > len(trackIDs) {
			upperLimit = len(trackIDs)
		}
		audioFeatures := getInfo.spotifyHandler.GetAudioFeaturesOfTrack(trackIDs[i:upperLimit])
		for _, a := range audioFeatures {
			if _, ok := getInfo.audioFeatures[a.ID]; !ok {
				getInfo.audioFeatures[a.ID] = a
			}
		}
		if i != 0 && i%500 == 0 {
			fmt.Println("5 second audio feature sleep")
			fmt.Println("AudioFeatures index: ", i)
			time.Sleep(5 * time.Second)
		}
		i += 50
	}
	if _, ok := getInfo.audioFeatures[constants.NotFound]; !ok {
		getInfo.audioFeatures[constants.NotFound] = nil
	}

	fmt.Println("-----------------------------")
	fmt.Println("Total unique songs found: ", len(trackIDs))
}

func (getInfo *getRecentTrackInformation) printoutResults() {
	f, err := os.Create("audioFeaturesTimeSeries_" + constants.Now.Format("20060102150405") + ".txt")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	dataWriter := bufio.NewWriter(f)

	_, _ = dataWriter.WriteString("ListenDate\tTrack\tAlbum\tArtist\tDuration(S)\tSpotifyID\tAcousticness\tDanceability\tEnergy\tInstrumentalness\tLiveness\tLoudness\tSpeechiness\tTempo\tValence\n")
	for _, t := range getInfo.tracksForDuration {
		af := getInfo.audioFeatures[t.SpotifyID]
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
		if af == nil{
			for i := 0; i < 10; i++ {
				trackStringArray = append(trackStringArray, constants.DoubleHyphen)
			}
		} else {
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
		}

		_, _ = dataWriter.WriteString(strings.Join(trackStringArray[:], "\t") + "\n")
	}

	_ = dataWriter.Flush()
	_ = f.Close()
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
