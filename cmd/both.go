package cmd

import "github.com/spf13/cobra"

func NewBothCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "both",
		Short: "Runs subcommands based on LastFM and Spotify data.",
	}
}
