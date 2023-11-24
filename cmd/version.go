package cmd

import (
	internal "diagram-converter/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commit         = "nnnnnnn"
	gitHubRepoUser = "sindrel"
	gitHubRepoName = "excalidraw-converter"
	noVersionCheck bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the application version",
	Long:  `This command provides information about the release version and the Git commit it was built from.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s (%s)\n", version, string(commit[0:7]))

		if !noVersionCheck {
			internal.PrintVersionCheck(gitHubRepoUser, gitHubRepoName, version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&noVersionCheck, "no-check", false, "opt-out from checking for latest version")
}
