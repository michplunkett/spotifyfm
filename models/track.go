package models

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shkh/lastfm-go/lastfm"
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/util/constants"
)

const (
	songAudioFeaturesFileName = "spotifyIDToAudioFeature.json"
	songListFileName          = "lastFMTrackListing.json"
	songIDsFileName           = "spotifySearchStringToSongID.json"
)

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
	SpotifyID       spotify.ID
}

func UserGetRecentTracksToDomainTracks(trackList *lastfm.UserGetRecentTracks) []Track {
	tracks := make([]Track, 0)
	for _, lastFMTrack := range trackList.Tracks {
		var listenDate time.Time
		if lastFMTrack.Date.Uts != constants.EmptyString {
			utcInt64, err := strconv.ParseInt(lastFMTrack.Date.Uts, 10, 64)
			if err != nil {
				panic(err)
			}
			listenDate = time.Unix(utcInt64, 0)
		} else {
			listenDate = time.Now()
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
			duration = constants.AverageSongSeconds
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

func GetSpotifySearchToSongIDs() map[string]spotify.ID {
	searchToIDHash := make(map[string]spotify.ID, 0)

	songIDFile, err := ioutil.ReadFile(songIDsFileName)
	if err != nil {
		return searchToIDHash
	}

	_ = json.Unmarshal(songIDFile, &searchToIDHash)

	return searchToIDHash
}

func AddSpotifySearchToSongIDs(searchToID map[string]spotify.ID) {
	validSearches := make(map[string]spotify.ID, 0)
	for key, id := range searchToID {
		if _, ok := validSearches[key]; !ok && id != constants.NotFound {
			validSearches[key] = id
		}
	}

	file, _ := json.MarshalIndent(validSearches, "", " ")
	_ = ioutil.WriteFile(songIDsFileName, file, 0644)
}

func GetSpotifyIDToAudioFeatures() map[spotify.ID]*spotify.AudioFeatures {
	audioFeatures := make(map[spotify.ID]*spotify.AudioFeatures, 0)

	songIDFile, err := ioutil.ReadFile(songAudioFeaturesFileName)
	if err != nil {
		return audioFeatures
	}

	_ = json.Unmarshal(songIDFile, &audioFeatures)

	return audioFeatures
}

func AddSpotifyIDToAudioFeatures(idToAudioFeatures map[spotify.ID]*spotify.AudioFeatures) {
	file, _ := json.MarshalIndent(idToAudioFeatures, "", " ")
	_ = ioutil.WriteFile(songAudioFeaturesFileName, file, 0644)
}

func AddLastFMTrackList(tracks []Track) {
	file, _ := json.MarshalIndent(tracks, "", " ")
	_ = ioutil.WriteFile(songListFileName, file, 0644)
}
