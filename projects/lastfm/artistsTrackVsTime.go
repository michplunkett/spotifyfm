package lastfm

import (
	"fmt"
	"sort"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

// There are some artists that have aliases in last.fm
var problematicArtistMapping = map[string]string{
	"strfkr": "starfucker",
	"th1rt3en": "thirteen",
}

type ArtistsTrackVsTime interface {
	Execute()
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

func (a *artistsTrackVsTime) Execute() {
	a.getInformation()
	a.doCalculations()
	a.printoutResults()
}

func (a *artistsTrackVsTime) getInformation() {
	a.tracks = a.handler.GetAllTopTracks(constants.APIObjectLimit, a.timeSpan, a.userName)
	a.lastFMSortedArtists = a.handler.GetAllTopArtists(constants.APIObjectLimit, a.timeSpan, a.userName)

	// UUIDs are mapped to their respective artist's lower case name.
	artistNameToUUIDHash := make(map[string]string, 0)
	for _, artist := range a.lastFMSortedArtists {
		if _, ok := artistNameToUUIDHash[artist.LowerCaseName]; !ok && artist.UUID != constants.EmptyString {
			artistNameToUUIDHash[artist.LowerCaseName] = artist.UUID
		}
	}

	// Add artist UUIDs to Track structs that are missing their respective artist's UUID.
	for idx, track := range a.tracks {
		if track.ArtistUUID == constants.EmptyString {
			if uuid, ok := artistNameToUUIDHash[track.LowerCaseArtist]; ok {
				track.ArtistUUID = uuid
				a.tracks[idx] = track
			} else {
				if hashName, ok := problematicArtistMapping[track.LowerCaseArtist]; ok {
					if uuid, ok := artistNameToUUIDHash[hashName]; ok {
						track.ArtistUUID = uuid
						a.tracks[idx] = track
					}
				} else {
					// Print any track that does not have a present artist.
					fmt.Println(track)
				}
			}
		}
	}
}

func (a *artistsTrackVsTime) doCalculations() {
	// Create a hash mapping an artist's UUID to its track duration sum.
	artistsDurationHash := make(map[string]int, 0)
	for _, artist := range a.lastFMSortedArtists {
		artistsDurationHash[artist.UUID] = 0
	}

	// Add up the duration of tracks listened per artist.
	for _, track := range a.tracks {
		if artist, ok := artistsDurationHash[track.ArtistUUID]; ok {
			// The total time listened to one song is its duration * number of times played.
			artist += track.Duration * track.PlayCount
			artistsDurationHash[track.ArtistUUID] = artist
		}
	}

	// Duration sum values are updated in the Artist struct.
	for idx, artist := range a.lastFMSortedArtists {
		if durationSum, ok := artistsDurationHash[artist.UUID]; ok {
			artist.DurationSum = durationSum
			a.lastFMSortedArtists[idx] = artist
		}
	}

	// Copy the lastFMSortedArtists array to durationSortedArtists.
	a.durationSortedArtists = append(a.durationSortedArtists, a.lastFMSortedArtists...)
	// Sort durationSortedArtists in DESC order by the DurationSum value.
	sort.SliceStable(a.durationSortedArtists, func(i, j int) bool {
		return a.durationSortedArtists[i].DurationSum > a.durationSortedArtists[j].DurationSum
	})

}

func (a *artistsTrackVsTime) printoutResults() {
	// I'd like a better way for visualising this, but this will have to do for right now.
	fmt.Println("---TOP 20 ARTISTS BY TRACKS LISTENED VS TIME LISTENED---")
	fmt.Printf("Rank\tArtist(Tracks)\tTrack Listens\tArtist(Time)\tMinutes Listened\n")
	for i := 0; i < 20; i++ {
		fmt.Printf("%d\t%s\t%d\t%s\t%d\n", i+1, a.lastFMSortedArtists[i].Name, a.lastFMSortedArtists[i].PlayCount,
			a.durationSortedArtists[i].Name, a.durationSortedArtists[i].DurationSum/60)
	}
}
