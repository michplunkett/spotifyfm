package endpoints

import "github.com/zmb3/spotify"

type SpotifyCaller interface {
	GetUserInfo() *spotify.PrivateUser
	GetCurrentlyPlaying() *spotify.CurrentlyPlaying
	GetTopTracks(timeRange string, pageSize int) *spotify.FullTrackPage
	GetTopArtists(timeRange string, pageSize int) *spotify.FullArtistPage
}

type spotifyHandler struct {
	client *spotify.Client
}

func NewSpotifyHandler(client *spotify.Client) *spotifyHandler {
	return &spotifyHandler{
		client: client,
	}
}

func (handler *spotifyHandler) GetUserInfo() *spotify.PrivateUser {
	userInfo, _ := handler.client.CurrentUser()
	return userInfo
}

func (handler *spotifyHandler) GetCurrentlyPlaying() *spotify.CurrentlyPlaying {
	currentlyPlaying, _ := handler.client.PlayerCurrentlyPlaying()
	return currentlyPlaying
}

func (handler *spotifyHandler) GetTopTracks(timeRange string, pageSize int) *spotify.FullTrackPage {
	topTracks, _ := handler.client.CurrentUsersTopTracksOpt(&spotify.Options{
		Limit:     &pageSize,
		Timerange: &timeRange,
	})
	return topTracks
}

func (handler *spotifyHandler) GetTopArtists(timeRange string, pageSize int) *spotify.FullArtistPage {
	topArtists, _ := handler.client.CurrentUsersTopArtistsOpt(&spotify.Options{
		Limit:     &pageSize,
		Timerange: &timeRange,
	})
	return topArtists
}
