package endpoints

import "github.com/zmb3/spotify"

type SpotifyCaller interface {
	GetUserInfo() *spotify.PrivateUser
	GetTopTracks() *spotify.FullTrackPage
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

func (handler *spotifyHandler) GetTopTracks(timeRange string, pageSize int) *spotify.FullTrackPage {
	topTracks, _ := handler.client.CurrentUsersTopTracksOpt(&spotify.Options{
		Limit:     &pageSize,
		Timerange: &timeRange,
	})
	return topTracks
}
