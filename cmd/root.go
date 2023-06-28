package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exconv",
	Short: "Excalidraw diagram converter.",
	Long: `A command line tool for porting Excalidraw diagrams to other formats. 

Excalidraw Converter ports Excalidraw diagrams to other formats, 
such as the .gliffy format for Gliffy and Gliffy for Confluence.

See the project repository for more information:
https://github.com/sindrel/excalidraw-converter`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
