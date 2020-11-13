package endpoints

import (
	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMHandler interface {
	GetUserInfo() *lastfm.UserGetInfo
	GetCurrentTrack(userName string, limit int) *lastfm.UserGetRecentTracks
	GetTopTracks(userName string, limit int, period string) *lastfm.UserGetTopTracks
	GetTopArtists(userName string, limit int, period string) *lastfm.UserGetTopArtists
}

type lastFMHandler struct {
	api *lastfm.Api
}

func NewLastFMHandler(api *lastfm.Api) *lastFMHandler {
	return &lastFMHandler{
		api: api,
	}
}

func (handler *lastFMHandler) GetUserInfo() *lastfm.UserGetInfo {
	userInfo, _ := handler.api.User.GetInfo(nil)
	return &userInfo
}

func (handler *lastFMHandler) GetCurrentTrack(userName string, limit int) *lastfm.UserGetRecentTracks {
	currentTrackParam := make(map[string]interface{})
	currentTrackParam["user"] = userName
	currentTrackParam["limit"] = limit
	currentTrack, _ := handler.api.User.GetRecentTracks(currentTrackParam)
	return &currentTrack
}

func (handler *lastFMHandler) GetTopTracks(userName string, limit int, period string) *lastfm.UserGetTopTracks {
	topTracksParam := make(map[string]interface{})
	topTracksParam["user"] = userName
	topTracksParam["limit"] = limit
	topTracksParam["period"] = period
	topTracks, _ := handler.api.User.GetTopTracks(topTracksParam)
	return &topTracks
}

func (handler *lastFMHandler) GetTopArtists(userName string, limit int, period string) *lastfm.UserGetTopArtists {
	topArtistsParam := make(map[string]interface{})
	topArtistsParam["user"] = userName
	topArtistsParam["limit"] = limit
	topArtistsParam["period"] = period
	topArtists, _ := handler.api.User.GetTopArtists(topArtistsParam)
	return &topArtists
}
