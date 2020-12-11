package models

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/shkh/lastfm-go/lastfm"
)

type Artist struct {
	DurationSum   int
	LowerCaseName string
	Name          string
	PlayCount     int
	Rank          int
	UUID          string
}

func UserGetTopArtistsToDomainArtists(artistList *lastfm.UserGetTopArtists) []Artist {
	artists := make([]Artist, 0)
	for _, lastFMArtist := range artistList.Artists {
		playCount, _ := strconv.Atoi(lastFMArtist.PlayCount)
		rank, _ := strconv.Atoi(lastFMArtist.Rank)
		artist := Artist{
			Name:          lastFMArtist.Name,
			LowerCaseName: strings.ToLower(lastFMArtist.Name),
			PlayCount:     playCount,
			Rank:          rank,
			UUID:          lastFMArtist.Mbid,
		}
		// not all artists have mbID
		if artist.UUID == "" {
			artist.UUID = uuid.New().String()
		}
		artists = append(artists, artist)
	}
	return artists
}

type Track struct {
	Artist          string
	ArtistUUID      string
	Duration        int
	LowerCaseArtist string
	Name            string
	PlayCount       int
	Rank            int
}

func UserGetTopTracksToDomainTracks(trackList *lastfm.UserGetTopTracks) []Track {
	tracks := make([]Track, 0)
	for _, lastFMTrack := range trackList.Tracks {
		duration, _ := strconv.Atoi(lastFMTrack.Duration)
		playCount, _ := strconv.Atoi(lastFMTrack.PlayCount)
		rank, _ := strconv.Atoi(lastFMTrack.Rank)
		if duration == 0 {
			// Last.fm says the average length of a song is 3.5 minutes, 210 seconds.
			duration = 210
		}
		if playCount == 0 {
			// If the song is in the list, I presume it's been played at least once.
			playCount = 1
		}
		track := Track{
			Artist:          lastFMTrack.Artist.Name,
			ArtistUUID:      lastFMTrack.Artist.Mbid,
			Duration:        duration,
			LowerCaseArtist: strings.ToLower(lastFMTrack.Artist.Name),
			Name:            lastFMTrack.Name,
			PlayCount:       playCount,
			Rank:            rank,
		}
		tracks = append(tracks, track)
	}
	return tracks
}
