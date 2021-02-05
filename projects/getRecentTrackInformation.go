package projects

import (
	"fmt"
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
	audioFeatures     []*spotify.AudioFeatures
}

func NewGetRecentTrackInformation(fromDate int64, lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) GetRecentTrackInformation {
	return &getRecentTrackInformation{
		fromDate:          fromDate,
		lastFMHandler:     lastFMHandler,
		tracksForDuration: make([]models.Track, 0),
		spotifyHandler:    spotifyHandler,
		audioFeatures:     make([]*spotify.AudioFeatures, 0),
	}
}

func (getInfo *getRecentTrackInformation) Execute() {
	getInfo.getInformation()
}

func (getInfo *getRecentTrackInformation) getInformation() {
	getInfo.tracksForDuration = getInfo.lastFMHandler.GetAllRecentTracks(getInfo.fromDate, getInfo.lastFMHandler.GetUserInfo().Name)
	trackIDs := make([]spotify.ID, 0)
	couldNotFindInSearch := 0
	couldNotMatchInSearch := 0
	trackToIDHash := make(map[string]spotify.ID, 0)
	for i, track := range getInfo.tracksForDuration {
		if i%500 == 0 {
			fmt.Println("Sleepin' for 30 seconds so Spotify doesn't hate me.")
			time.Sleep(30 * time.Second)
		}

		searchKey := track.Artist + " " + track.AlbumName + " " + track.Name
		if value, ok := trackToIDHash[searchKey]; ok {
			if value != constants.NotFound {
				track.SpotifyID = value
				getInfo.tracksForDuration[i] = track
			} else {
				track.SpotifyID = constants.NotFound
				getInfo.tracksForDuration[i] = track
			}
			continue
		}

		searchResult := getInfo.spotifyHandler.SearchForSong(track.Artist, track.AlbumName, track.Name)
		if searchResult != nil {
			comparisonResult := compareMultipleReturnedTracks(track, searchResult)
			if comparisonResult != constants.EmptyString {
				trackIDs = append(trackIDs, comparisonResult)
				track.SpotifyID = comparisonResult
				getInfo.tracksForDuration[i] = track
				trackToIDHash[searchKey] = comparisonResult
			} else {
				fmt.Println("--- Count not match in search ---")
				fmt.Println("index: ", i)
				fmt.Println("search string: " + searchKey)
				fmt.Println(searchResult)
				track.SpotifyID = constants.NotFound
				trackToIDHash[searchKey] = constants.NotFound
				getInfo.tracksForDuration[i] = track
				couldNotMatchInSearch += 1
			}
		} else {
			fmt.Println("--- Count not find in search ---")
			fmt.Println("index: ", i)
			fmt.Println("search string: " + searchKey)
			track.SpotifyID = constants.NotFound
			trackToIDHash[searchKey] = constants.NotFound
			getInfo.tracksForDuration[i] = track
			couldNotFindInSearch += 1
		}
	}
	fmt.Println("Could not match in search: ", couldNotMatchInSearch)
	fmt.Println("Could not find in search: ", couldNotFindInSearch)
	fmt.Println("Total tracks: ", len(getInfo.tracksForDuration))
}

func (getInfo *getRecentTrackInformation) doCalculations() {

}

func (getInfo *getRecentTrackInformation) printoutResults() {

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
