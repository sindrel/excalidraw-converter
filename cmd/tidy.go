package cmd

import (
	conv "diagram-converter/internal/conversion"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var defaultOutputPathTidy = "your_file_0.excalidraw"

var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Snap the diagram to a grid",
	Long: `This command is used to tidy up a diagram by snapping it to a grid.

  Blabla.

Example:
  exconv tidy -i your_file.excalidraw
`,
	Run: func(cmd *cobra.Command, args []string) {
		importPath, _ := cmd.Flags().GetString("input")
		exportPath, _ := cmd.Flags().GetString("output")

		if len(importPath) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Input file path not provided.\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if strings.HasPrefix(exportPath, defaultOutputPath) {
			exportPath = strings.TrimSuffix(path.Base(importPath), filepath.Ext(importPath)) + "_0.excalidraw"
		}

		err := conv.ConvertExcalidrawToGliffyFile(importPath, exportPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to clean up Excalidraw diagram: %s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(gliffyCmd)

	gliffyCmd.PersistentFlags().StringP("input", "i", "", "input file path")
	gliffyCmd.PersistentFlags().StringP("output", "o", defaultOutputPathTidy, "output file path")
}
