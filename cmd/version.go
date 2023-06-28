package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "nnnnnnn"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the application version",
	Long:  `Provides information about the release version and the Git commit it was built from.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", version, string(commit[0:7]))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
