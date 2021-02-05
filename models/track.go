package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/util/constants"
)

type Track struct {
	AlbumName       string
	Artist          string
	ArtistUUID      string
	Duration        int
	ListenDate      int64
	LowerCaseArtist string
	Name            string
	PlayCount       int
	Rank            int
	SpotifyID       string
}

func UserGetRecentTracksToDomainTracks(trackList *lastfm.UserGetRecentTracks) []Track {
	tracks := make([]Track, 0)
	for _, lastFMTrack := range trackList.Tracks {
		var listenDate int64
		var err error
		if lastFMTrack.Date.Uts != constants.EmptyString {
			listenDate, err = strconv.ParseInt(lastFMTrack.Date.Uts, 10, 64)
			if err != nil {
				panic(err)
			}
		} else {
			listenDate = time.Now().Unix()
		}

		track := Track{
			AlbumName:       lastFMTrack.Album.Name,
			Artist:          lastFMTrack.Artist.Name,
			ArtistUUID:      lastFMTrack.Artist.Mbid,
			ListenDate:      listenDate,
			LowerCaseArtist: RemoveNonWordCharacters(lastFMTrack.Artist.Name),
			Name:            lastFMTrack.Name,
		}
		tracks = append(tracks, track)
	}
	return tracks
}

func UserGetTopTracksToDomainTracks(trackList *lastfm.UserGetTopTracks) []Track {
	tracks := make([]Track, 0)
	for _, lastFMTrack := range trackList.Tracks {
		duration, _ := strconv.Atoi(lastFMTrack.Duration)
		playCount, _ := strconv.Atoi(lastFMTrack.PlayCount)
		rank, _ := strconv.Atoi(lastFMTrack.Rank)
		// Not all songs come with assigned durations.
		if duration == 0 {
			// Last.fm says the average length of a song is 3.5 minutes, 210 seconds.
			duration = 210
		}
		// If the song is included, I presume it is on the list because it has been played AT LEAST once.
		if playCount == 0 {
			// If the song is in the list, I presume it's been played at least once.
			playCount = 1
		}
		track := Track{
			Artist:          lastFMTrack.Artist.Name,
			ArtistUUID:      lastFMTrack.Artist.Mbid,
			Duration:        duration,
			LowerCaseArtist: RemoveNonWordCharacters(lastFMTrack.Artist.Name),
			Name:            lastFMTrack.Name,
			PlayCount:       playCount,
			Rank:            rank,
		}
		tracks = append(tracks, track)
	}
	return tracks
}

var regEx, _ = regexp.Compile("[^a-zA-Z0-9]+")

func RemoveNonWordCharacters(name string) string {
	return regEx.ReplaceAllString(strings.ToLower(name), constants.EmptyString)
}
