package endpoints

import (
	"github.com/zmb3/spotify"
)

const (
	pageSizeConst int = 20
)

type SpotifyHandler interface {
	GetUserInfo() *spotify.PrivateUser
	GetCurrentlyPlaying() *spotify.CurrentlyPlaying
	GetAllTopTracks(timeRange string) []spotify.FullTrack
	GetAllTopArtists(timeRange string) []spotify.FullArtist
	GetAudioFeaturesOfTrack(trackIDs []spotify.ID) []*spotify.AudioFeatures
	SearchForSong(songName, albumName, artistName string) ([]spotify.FullTrack, error)
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

func (handler *spotifyHandler) getTopTracks(timeRange string, offset, pageSize int) []spotify.FullTrack {
	topTracks, _ := handler.client.CurrentUsersTopTracksOpt(&spotify.Options{
		Limit:     &pageSize,
		Offset:    &offset,
		Timerange: &timeRange,
	})
	return topTracks.Tracks
}

func (handler *spotifyHandler) GetAllTopTracks(timeRange string) []spotify.FullTrack {
	tracks := make([]spotify.FullTrack, 0)

	offset := 0
	for {
		topTracks := handler.getTopTracks(timeRange, offset, pageSizeConst)
		tracks = append(tracks, topTracks...)
		// When the amount of tracks being returned is less than the limit there are no more tracks to pull
		if len(topTracks) < pageSizeConst {
			break
		}
		offset = offset + pageSizeConst
	}
	return tracks
}

func (handler *spotifyHandler) getTopArtists(timeRange string, pageSize int) []spotify.FullArtist {
	topArtists, _ := handler.client.CurrentUsersTopArtistsOpt(&spotify.Options{
		Limit:     &pageSize,
		Timerange: &timeRange,
	})
	return topArtists.Artists
}

func (handler *spotifyHandler) GetAllTopArtists(timeRange string) []spotify.FullArtist {
	// There currently isn't an offset option for CurrentUsersTopArtistsOpt so I'm doing one large grab
	return handler.getTopArtists(timeRange, pageSizeConst)
}

func (handler *spotifyHandler) GetAudioFeaturesOfTrack(trackIDs []spotify.ID) []*spotify.AudioFeatures {
	audioFeatures, _ := handler.client.GetAudioFeatures(trackIDs...)
	if len(audioFeatures) > 0 {
		return audioFeatures
	}
	return nil
}

func (handler *spotifyHandler) SearchForSong(artistName, albumName, songName string) ([]spotify.FullTrack, error) {
	queryString := artistName + " " + albumName + " " + songName
	limit := 3
	options := spotify.Options{
		Limit: &limit,
	}
	result, err := handler.client.SearchOpt(queryString, spotify.SearchTypeTrack, &options)
	if result != nil && len(result.Tracks.Tracks) > 0 {
		return result.Tracks.Tracks, nil
	}
	return nil, err
}
