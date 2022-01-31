package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "search",
		Short: "Search messages by date",
		Long: `Search recently sent messages and optionally narrow by date range, tags, senders, and API keys.
		If no date range is specified, results within the last 7 days are returned.
		This method may be called up to 20 times per minute.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
