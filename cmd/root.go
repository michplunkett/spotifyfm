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
	var rootCmd = &cobra.Command{
		Use:     "spotifyfm",
		Version: "0.0.1",
		Short:   "SpotifyFM is a CLI application designed to compare Spotify and Last.fm calculations",
	}

	for cmd, subcommands := range map[*cobra.Command][]*cobra.Command{
		NewBothCmd(): {
			NewRecentTrackInformationCmd(lastFMHandler, spotifyHandler),
		},
		NewLastFMCmd(): {
			NewTrackVsTimeCmd(lastFMHandler),
		},
		NewSpotifyCmd(): {
			NewAudioFeatureProcessing(spotifyHandler),
		},
	} {
		// add subcommands to their parent
		for _, subCmd := range subcommands {
			cmd.AddCommand(subCmd)
		}

		// add parent command to the root
		rootCmd.AddCommand(cmd)
	}

	return &Root{cmd: rootCmd}
}

func (r *Root) Execute() {
	if err := r.cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
