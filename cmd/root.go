package cmd

import (
	"fmt"
	"os"

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
		for _, subCmd := range subcommands {
			cmd.AddCommand(subCmd)
		}

		// add parent command to the root
		root.AddCommand(cmd)
	}

	return &Root{cmd: root}
}

func (r *Root) Execute() {
	if err := r.cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
