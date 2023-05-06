package endpoints

import (
	"fmt"

	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

type LastFMHandler interface {
	GetAllRecentTracks(from int64, userName string) []models.Track
	GetAllTopArtists(period, userName string) []models.Artist
	GetAllTopTracks(period, userName string) []models.Track
	GetCurrentTrack(userName string) *lastfm.UserGetRecentTracks
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

func (handler *lastFMHandler) GetAllRecentTracks(from int64, userName string) []models.Track {
	tracks := make([]models.Track, 0)

	totalTracksFetched := 0
	page := 1
	for {
		topTracks := handler.getRecentTracks(from, constants.APIObjectLimit, page, userName)
		domainTracks := models.UserGetRecentTracksToDomainTracks(topTracks)
		tracks = append(tracks, domainTracks...)
		// When the amount of tracks being returned is less than the limit there are no more tracks to pull
		if len(domainTracks) < constants.APIObjectLimit {
			break
		}
		page += 1
		totalTracksFetched += len(domainTracks)
		if totalTracksFetched%1000 == 0 {
			fmt.Printf("You have fetched %d tracks from Last.FM.\n", totalTracksFetched)
		}
	}
	return tracks
}

func (handler *lastFMHandler) GetAllTopArtists(period, userName string) []models.Artist {
	artists := make([]models.Artist, 0)

	page := 1
	for {
		topArtists := handler.getTopArtists(constants.APIObjectLimit, page, period, userName)
		domainArtists := models.UserGetTopArtistsToDomainArtists(topArtists)
		artists = append(artists, domainArtists...)
		// When the amount of artists being returned is less than the limit there are no more artists to pull
		if len(domainArtists) < constants.APIObjectLimit {
			break
		}
		page = page + 1
	}
	return artists
}

func (handler *lastFMHandler) GetAllTopTracks(period, userName string) []models.Track {
	tracks := make([]models.Track, 0)

	page := 1
	for {
		topTracks := handler.getTopTracks(constants.APIObjectLimit, page, period, userName)
		domainTracks := models.UserGetTopTracksToDomainTracks(topTracks)
		tracks = append(tracks, domainTracks...)
		// When the amount of tracks being returned is less than the limit there are no more tracks to pull
		if len(domainTracks) < constants.APIObjectLimit {
			break
		}
		page = page + 1
	}
	return tracks
}

func (handler *lastFMHandler) GetCurrentTrack(userName string) *lastfm.UserGetRecentTracks {
	currentTrackParam := make(map[string]interface{})
	currentTrackParam["user"] = userName
	currentTrackParam["limit"] = constants.APIObjectLimit
	currentTrack, _ := handler.api.User.GetRecentTracks(currentTrackParam)
	return &currentTrack
}

func (handler *lastFMHandler) getRecentTracks(from int64, limit, page int, userName string) *lastfm.UserGetRecentTracks {
	topArtistsParam := make(map[string]interface{})
	topArtistsParam["from"] = from
	topArtistsParam["limit"] = limit
	topArtistsParam["page"] = page
	topArtistsParam["user"] = userName
	recentTracks, _ := handler.api.User.GetRecentTracks(topArtistsParam)
	return &recentTracks
}

func (handler *lastFMHandler) getTopArtists(limit, page int, period, userName string) *lastfm.UserGetTopArtists {
	topArtistsParam := make(map[string]interface{})
	topArtistsParam["limit"] = limit
	topArtistsParam["page"] = page
	topArtistsParam["period"] = period
	topArtistsParam["user"] = userName
	topArtists, _ := handler.api.User.GetTopArtists(topArtistsParam)
	return &topArtists
}

func (handler *lastFMHandler) getTopTracks(limit, page int, period string, userName string) *lastfm.UserGetTopTracks {
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
