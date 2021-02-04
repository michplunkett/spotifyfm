package endpoints

import (
	"time"

	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/models"
)

type LastFMHandler interface {
	GetAllRecentTracks(from time.Time, limit int, userName string) []models.Track
	GetAllTopArtists(limit int, period string, userName string) []models.Artist
	GetAllTopTracks(limit int, period string, userName string) []models.Track
	GetCurrentTrack(userName string, limit int) *lastfm.UserGetRecentTracks
	GetRecentTracks(from time.Time, limit, page int, userName string) *lastfm.UserGetRecentTracks
	GetTopArtists(limit, page int, period, userName string) *lastfm.UserGetTopArtists
	GetTopTracks(limit, page int, period string, userName string) *lastfm.UserGetTopTracks
	GetUserInfo() *lastfm.UserGetInfo
}

type lastFMHandler struct {
	api *lastfm.Api
}

func NewLastFMHandler(api *lastfm.Api) LastFMHandler {
	return &lastFMHandler{
		api: api,
	}
}

func (handler *lastFMHandler) GetAllRecentTracks(from time.Time, limit int, userName string) []models.Track {
	tracks := make([]models.Track, 0)

	page := 1
	for {
		topTracks := handler.GetRecentTracks(from, limit, page, userName)
		domainTracks := models.UserGetRecentTracksToDomainTracks(topTracks)
		tracks = append(tracks, domainTracks...)
		// When the amount of tracks being returned is less than the limit there are no more tracks to pull
		if len(domainTracks) < limit {
			break
		}
		page = page + 1
	}
	return tracks
}

func (handler *lastFMHandler) GetAllTopArtists(limit int, period string, userName string) []models.Artist {
	artists := make([]models.Artist, 0)

	page := 1
	for {
		topArtists := handler.GetTopArtists(limit, page, period, userName)
		domainArtists := models.UserGetTopArtistsToDomainArtists(topArtists)
		artists = append(artists, domainArtists...)
		// When the amount of artists being returned is less than the limit there are no more artists to pull
		if len(domainArtists) < limit {
			break
		}
		page = page + 1
	}
	return artists
}

func (handler *lastFMHandler) GetAllTopTracks(limit int, period string, userName string) []models.Track {
	tracks := make([]models.Track, 0)

	page := 1
	for {
		topTracks := handler.GetTopTracks(limit, page, period, userName)
		domainTracks := models.UserGetTopTracksToDomainTracks(topTracks)
		tracks = append(tracks, domainTracks...)
		// When the amount of tracks being returned is less than the limit there are no more tracks to pull
		if len(domainTracks) < limit {
			break
		}
		page = page + 1
	}
	return tracks
}

func (handler *lastFMHandler) GetCurrentTrack(userName string, limit int) *lastfm.UserGetRecentTracks {
	currentTrackParam := make(map[string]interface{})
	currentTrackParam["user"] = userName
	currentTrackParam["limit"] = limit
	currentTrack, _ := handler.api.User.GetRecentTracks(currentTrackParam)
	return &currentTrack
}

func (handler *lastFMHandler) GetRecentTracks(from time.Time, limit, page int, userName string) *lastfm.UserGetRecentTracks {
	topArtistsParam := make(map[string]interface{})
	topArtistsParam["from"] = from
	topArtistsParam["limit"] = limit
	topArtistsParam["page"] = page
	topArtistsParam["user"] = userName
	recentTracks, _ := handler.api.User.GetRecentTracks(topArtistsParam)
	return &recentTracks
}

func (handler *lastFMHandler) GetTopArtists(limit, page int, period, userName string) *lastfm.UserGetTopArtists {
	topArtistsParam := make(map[string]interface{})
	topArtistsParam["limit"] = limit
	topArtistsParam["page"] = page
	topArtistsParam["period"] = period
	topArtistsParam["user"] = userName
	topArtists, _ := handler.api.User.GetTopArtists(topArtistsParam)
	return &topArtists
}

func (handler *lastFMHandler) GetTopTracks(limit, page int, period string, userName string) *lastfm.UserGetTopTracks {
	topTracksParam := make(map[string]interface{})
	topTracksParam["limit"] = limit
	topTracksParam["page"] = page
	topTracksParam["period"] = period
	topTracksParam["user"] = userName
	topTracks, _ := handler.api.User.GetTopTracks(topTracksParam)
	return &topTracks
}

func (handler *lastFMHandler) GetUserInfo() *lastfm.UserGetInfo {
	userInfo, _ := handler.api.User.GetInfo(nil)
	return &userInfo
}
