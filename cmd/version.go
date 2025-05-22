package cmd

import (
	internal "diagram-converter/internal"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commit         = ""
	githubRepoUser = "sindrel"
	githubRepoName = "excalidraw-converter"
	noVersionCheck bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the application version",
	Long:  `This command provides information about current and available release(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version)
		if commit != "" {
			fmt.Printf("%s\n", string(commit[0:7]))
		}

		if !noVersionCheck {
			err := internal.PrintVersionCheck(githubRepoUser, githubRepoName, "v"+version)
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
