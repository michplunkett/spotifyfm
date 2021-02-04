package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/util/constants"
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
			LowerCaseName: removeNonWordCharacters(strings.ToLower(lastFMArtist.Name)),
			PlayCount:     playCount,
			Rank:          rank,
			UUID:          lastFMArtist.Mbid,
		}
		// Not all artists have a UUID.
		if artist.UUID == "" {
			artist.UUID = uuid.New().String()
		}
		artists = append(artists, artist)
	}
	return artists
}

type Track struct {
	AlbumName       string
	Artist          string
	ArtistUUID      string
	Duration        int
	ListenDate      time.Time
	LowerCaseArtist string
	Name            string
	PlayCount       int
	Rank            int
	SpotifyID       string
}

func UserGetRecentTracksToDomainTracks(trackList *lastfm.UserGetRecentTracks) []Track {
	tracks := make([]Track, 0)
	for _, lastFMTrack := range trackList.Tracks {
		i, err := strconv.ParseInt(lastFMTrack.Date.Uts, 10, 64)
		if err != nil {
			panic(err)
		}
		listenDate := time.Unix(i, 0)

		track := Track{
			AlbumName:       lastFMTrack.Album.Name,
			Artist:          lastFMTrack.Artist.Name,
			ArtistUUID:      lastFMTrack.Artist.Mbid,
			ListenDate:      listenDate,
			LowerCaseArtist: removeNonWordCharacters(strings.ToLower(lastFMTrack.Artist.Name)),
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
			LowerCaseArtist: removeNonWordCharacters(strings.ToLower(lastFMTrack.Artist.Name)),
			Name:            lastFMTrack.Name,
			PlayCount:       playCount,
			Rank:            rank,
		}
		tracks = append(tracks, track)
	}
	return tracks
}

var regEx, _ = regexp.Compile("[^a-zA-Z0-9]+")

func removeNonWordCharacters(name string) string {
	return regEx.ReplaceAllString(name, constants.EmptyString)
}
