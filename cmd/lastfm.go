package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/michplunkett/spotifyfm/util/constants"
)

var (
	durationSort []models.Artist
	lastFMSort   []models.Artist
	timeSpan     = constants.LastFMPeriod6Month
	tracks       []models.Track
)

func NewLastFMCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lastfm",
		Short: "Runs LastFM subcommands.",
	}
}

// TODO: Needs flags etc.

func NewTrackVsTimeCmd(handler endpoints.LastFMHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "track-vs-time",
		Short: "Compares most common artist ranks by number of tracks versus overall time.",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting the comparison process...")
			getInformation(handler, timeSpan)
			fmt.Println(durationSort)
			fmt.Println(lastFMSort)
			doCalculations()
			printoutResults()
		},
	}

	return cmd
}

func getInformation(handler endpoints.LastFMHandler, timeSpan string) {
	userName := handler.GetUserInfo().Name
	tracks = handler.GetAllTopTracks(timeSpan, userName)
	lastFMSort = handler.GetAllTopArtists(timeSpan, userName)

	// UUIDs are mapped to their respective artist's lower case name.
	artistNameToUUIDHash := make(map[string]string, 0)
	for _, artist := range lastFMSort {
		if _, ok := artistNameToUUIDHash[artist.LowerCaseName]; !ok && artist.UUID != constants.EmptyString {
			artistNameToUUIDHash[artist.LowerCaseName] = artist.UUID
		}
	}

	// Add artist UUIDs to Track structs that are missing their respective artist's UUID.
	for idx, track := range tracks {
		if track.ArtistUUID == constants.EmptyString {
			if uuid, ok := artistNameToUUIDHash[track.LowerCaseArtist]; ok {
				track.ArtistUUID = uuid
				tracks[idx] = track
			} else {
				if hashName, problematicOk := constants.ProblematicArtistMapping[track.LowerCaseArtist]; problematicOk {
					if uuid, ok = artistNameToUUIDHash[hashName]; ok {
						track.ArtistUUID = uuid
						tracks[idx] = track
					}
				} else {
					// Print any track that does not have a present artist.
					fmt.Println(track)
				}
			}
		}
	}
}

func doCalculations() {
	// Create a hash mapping an artist's UUID to its track duration sum.
	artistsDurationHash := make(map[string]int, 0)
	for _, artist := range lastFMSort {
		artistsDurationHash[artist.UUID] = 0
	}

	// Add up the duration of tracks listened per artist.
	for _, track := range tracks {
		if artist, ok := artistsDurationHash[track.ArtistUUID]; ok {
			// The total time listened to one song is its duration * number of times played.
			artist += track.Duration * track.PlayCount
			artistsDurationHash[track.ArtistUUID] = artist
		}
	}

	// Duration sum values are updated in the Artist struct.
	for idx, artist := range lastFMSort {
		if durationSum, ok := artistsDurationHash[artist.UUID]; ok {
			artist.DurationSum = durationSum
			lastFMSort[idx] = artist
		}
	}

	// Copy the lastFMSortedArtists array to durationSortedArtists.
	durationSort = append(durationSort, lastFMSort...)
	// Sort durationSortedArtists in DESC order by the DurationSum value.
	sort.SliceStable(durationSort, func(i, j int) bool {
		return durationSort[i].DurationSum > durationSort[j].DurationSum
	})

}

func printoutResults() {
	// I'd like a better way for visualising this, but this will have to do for right now.
	fmt.Println("---TOP 20 ARTISTS BY TRACKS LISTENED VS TIME LISTENED---")
	fmt.Printf("Rank\tArtist(Tracks)\tTrack Listens\tArtist(Time)\tMinutes Listened\n")
	for i := 0; i < 20; i++ {
		fmt.Printf("%d\t%s\t%d\t%s\t%d\n", i+1, lastFMSort[i].Name, lastFMSort[i].PlayCount,
			durationSort[i].Name, durationSort[i].DurationSum/60)
	}
}
