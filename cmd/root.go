package cmd

import (
	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/spf13/cobra"
)

//projects.NewGetRecentTrackInformation(constants.StartOf2021, lastFMHandler, spotifyHandler).Execute()
//projects.NewAudioFeatureProcessing("audioFeaturesTimeSeries_20211014115400.txt", spotifyHandler).Execute()

type Root struct {
	cmd *cobra.Command
}

func NewRootCmd(lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) *Root {
	root := &cobra.Command{
		Use:   "spotifyfm",
		Short: "SpotifyFM is a CLI application designed to compare Spotify and Last.fm calculations",
	}

	for cmd, subcommands := range map[*cobra.Command][]*cobra.Command{
		NewBothCmd(): nil,
		NewLastFMCmd(): {
			NewTrackVsTimeCmd(lastFMHandler),
		},
		NewSpotifymd(): nil,
	} {
		// add subcommands to their parent
		for _, subcmd := range subcommands {
			cmd.AddCommand(subcmd)
		}

		// add parent command to the root
		root.AddCommand(cmd)
	}

	return &Root{cmd: root}
}
