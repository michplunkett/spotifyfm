package lastfm

import (
	"fmt"
	"sort"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

type ArtistsTrackVsTime interface {
	GetInformation()
	DoCalculations()
	PrintoutResults()
}

type artistsTrackVsTime struct {
	durationSortedArtists []models.Artist
	handler               endpoints.LastFMHandler
	lastFMSortedArtists   []models.Artist
	timeSpan              string
	tracks                []models.Track
	userName              string
}

func NewArtistsTrackVsTime(handler endpoints.LastFMHandler, timeSpan, userName string) ArtistsTrackVsTime {
	return &artistsTrackVsTime{
		durationSortedArtists: make([]models.Artist, 0),
		handler:               handler,
		lastFMSortedArtists:   make([]models.Artist, 0),
		timeSpan:              timeSpan,
		tracks:                make([]models.Track, 0),
		userName:              userName,
	}
}

func (a *artistsTrackVsTime) GetInformation() {
	a.tracks = a.handler.GetAllTopTracks(constants.APIObjectLimit, a.timeSpan, a.userName)
	a.lastFMSortedArtists = a.handler.GetAllTopArtists(constants.APIObjectLimit, a.timeSpan, a.userName)

	artistNameToUUIDHash := make(map[string]string, 0)
	for _, artist := range a.lastFMSortedArtists {
		if _, ok := artistNameToUUIDHash[artist.LowerCaseName]; !ok && artist.UUID != constants.EmptyString {
			artistNameToUUIDHash[artist.LowerCaseName] = artist.UUID
		}
	}
	for _, track := range a.tracks {
		if _, ok := artistNameToUUIDHash[track.LowerCaseArtist]; !ok && track.ArtistUUID != constants.EmptyString {
			artistNameToUUIDHash[track.LowerCaseArtist] = track.ArtistUUID
		}
	}

	for idx, track := range a.tracks {
		if track.ArtistUUID == constants.EmptyString {
			if mbID, ok := artistNameToUUIDHash[track.LowerCaseArtist]; ok {
				track.ArtistUUID = mbID
				a.tracks[idx] = track
			} else {
				fmt.Println(track)
			}
		}
	}
}

func (a *artistsTrackVsTime) DoCalculations() {
	artistsDurationHash := make(map[string]int, 0)
	for _, artist := range a.lastFMSortedArtists {
		artistsDurationHash[artist.UUID] = 0
	}

	for _, track := range a.tracks {
		if artist, ok := artistsDurationHash[track.ArtistUUID]; ok {
			artist += track.Duration * track.PlayCount
			artistsDurationHash[track.ArtistUUID] = artist
		} else {
			fmt.Println(track)
		}
	}

	for idx, artist := range a.lastFMSortedArtists {
		if durationSum, ok := artistsDurationHash[artist.UUID]; ok {
			artist.DurationSum = durationSum
			a.lastFMSortedArtists[idx] = artist
		}
	}

	a.durationSortedArtists = append(a.durationSortedArtists, a.lastFMSortedArtists...)
	sort.SliceStable(a.durationSortedArtists, func(i, j int) bool {
		return a.durationSortedArtists[i].DurationSum > a.durationSortedArtists[j].DurationSum
	})

}

func (a *artistsTrackVsTime) PrintoutResults() {
	fmt.Println("---TOP 20 ARTISTS BY TRACKS LISTENED VS TIME LISTENED---")
	fmt.Printf("Rank\tArtist(Tracks)\tTrack Listens\tArtist(Time)\tMinutes Listened\n")
	for i := 0; i < 20; i++ {
		fmt.Printf("%d\t%s\t%d\t%s\t%d\n", i+1, a.lastFMSortedArtists[i].Name, a.lastFMSortedArtists[i].PlayCount,
			a.durationSortedArtists[i].Name, a.durationSortedArtists[i].DurationSum/60)
	}
}
