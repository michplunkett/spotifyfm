package endpoints

import (
	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMHandler interface {
	GetUserInfo()
	GetTopTracks()
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

func (handler *lastFMHandler) GetTopTracks(userName string, limit int, period string) *lastfm.UserGetTopTracks {
	topTracksParam := make(map[string]interface{})
	topTracksParam["user"] = userName
	topTracksParam["limit"] = limit
	topTracksParam["period"] = period
	topTracks, _ := handler.api.User.GetTopTracks(topTracksParam)
	return &topTracks
}
