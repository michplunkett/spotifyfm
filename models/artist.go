package models

import (
	"strconv"

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
			LowerCaseName: RemoveNonWordCharacters(lastFMArtist.Name),
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
