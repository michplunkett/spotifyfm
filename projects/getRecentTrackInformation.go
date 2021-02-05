package projects

import (
	"fmt"
	"github.com/zmb3/spotify"
	"strings"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
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
	for i, track := range getInfo.tracksForDuration {
		searchResult := getInfo.spotifyHandler.SearchForSong(track.Artist, track.AlbumName, track.Name)
		if searchResult != nil {
			if strings.ToLower(searchResult.Name) == strings.ToLower(track.Name) &&
				models.RemoveNonWordCharacters(searchResult.Artists[0].Name) == track.LowerCaseArtist {
				trackIDs = append(trackIDs, searchResult.ID)
			} else {
				fmt.Println("--- Count not match in search ---")
				fmt.Println("index: ", i)
				fmt.Println("search string: " + track.LowerCaseArtist + " " + track.AlbumName + " " + track.Name)
				fmt.Println(models.RemoveNonWordCharacters(searchResult.Artists[0].Name), searchResult.Album.Name, searchResult.Name)
				couldNotMatchInSearch += 1
			}
		} else {
			fmt.Println("--- Count not find in search ---")
			fmt.Println("index: ", i)
			fmt.Println("search string: " + track.Artist + " " + track.AlbumName + " " + track.Name)
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
