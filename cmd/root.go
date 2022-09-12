package cmd

import (
	"github.com/spf13/cobra"

	"github.com/michplunkett/spotifyfm/api/endpoints"
)

type Root struct {
	cmd *cobra.Command
}

func NewRootCmd(lastFMHandler endpoints.LastFMHandler, spotifyHandler endpoints.SpotifyHandler) *Root {
	root := &cobra.Command{
		Use:   "spotifyfm",
		Short: "SpotifyFM is a CLI application designed to compare Spotify and Last.fm calculations",
	}

	for cmd, subcommands := range map[*cobra.Command][]*cobra.Command{
		NewBothCmd(): {
			NewRecentTrackInformationCmd(lastFMHandler, spotifyHandler),
		},
		NewLastFMCmd(): {
			NewTrackVsTimeCmd(lastFMHandler),
		},
		NewSpotifymd(): {
			NewAudioFeatureProcessing(spotifyHandler),
		},
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
