package cmd

import "github.com/spf13/cobra"

func NewSpotifymd() *cobra.Command {
	return &cobra.Command{
		Use:   "spotify",
		Short: "Runs Spotify subcommands.",
	}
}
