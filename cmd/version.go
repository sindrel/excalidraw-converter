package cmd

import (
	internal "diagram-converter/internal"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commit         = "nnnnnnn"
	githubRepoUser = "sindrel"
	githubRepoName = "excalidraw-converter"
	noVersionCheck bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the application version",
	Long:  `This command provides information about the release version and the Git commit it was built from.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", version, string(commit[0:7]))

		if !noVersionCheck {
			err := internal.PrintVersionCheck(githubRepoUser, githubRepoName, version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to check for latest version: %s\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&noVersionCheck, "no-check", false, "opt-out from checking for latest version")
}
