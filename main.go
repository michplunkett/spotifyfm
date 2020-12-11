package main

import (
	"fmt"
	"github.com/michplunkett/spotifyfm/api"
	"github.com/michplunkett/spotifyfm/api/authentication"
	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/config"
	"github.com/michplunkett/spotifyfm/projects/lastfm"
)

func main() {
	fmt.Println("---------------------")
	fmt.Println("Doin' some music things")

	// Starting the http server
	api.Start()

	// Getting the environment variables
	envVars := config.NewEnvVars()

	spotifyAuth := authentication.NewSpotifyAuthHandlerAll(envVars)
	spotifyClient := spotifyAuth.Authenticate()
	spotifyHandler := endpoints.NewSpotifyHandler(spotifyClient)

	fmt.Println("---------------------")
	fmt.Println("SPOTIFY THANGS")

	// use the client to make calls that require authorization
	spotifyUser := spotifyHandler.GetUserInfo()
	fmt.Println("You are logged in as ", spotifyUser.DisplayName)

	spotifyCurrentlyPlaying := spotifyHandler.GetCurrentlyPlaying()
	if spotifyCurrentlyPlaying.Item != nil {
		fmt.Println("This is the track you're currently playing ", spotifyCurrentlyPlaying.Item.Name)
	}

	spotifyTopTracks := spotifyHandler.GetTopTracks(config.SpotifyPeriodShort, 50)
	fmt.Println("This is your top track ", spotifyTopTracks.Tracks[0].SimpleTrack.Name)

	spotifyTopArtists := spotifyHandler.GetTopArtists(config.SpotifyPeriodShort, 50)
	fmt.Println("This is your top artist ", spotifyTopArtists.Artists[0].Name)

	lastFMAuth := authentication.NewLastFMAuthHandler(envVars)
	lastFMClient := lastFMAuth.Authenticate()
	lastFMHandler := endpoints.NewLastFMHandler(lastFMClient)

	fmt.Println("---------------------")
	fmt.Println("LAST.FM THANGS")

	lastFMUser := lastFMHandler.GetUserInfo()
	fmt.Println("You are logged in as ", lastFMUser.RealName)

	lastFMArtistsTrackVsTime := lastfm.NewArtistsTrackVsTime(lastFMHandler, config.LastFMPeriod1Month, lastFMUser.Name)
	lastFMArtistsTrackVsTime.GetInformation()
	lastFMArtistsTrackVsTime.DoCalculations()
	lastFMArtistsTrackVsTime.PrintoutResults()
}
