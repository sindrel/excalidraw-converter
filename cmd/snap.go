package cmd

import (
	snap "diagram-converter/internal/snapping"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var defaultOutputPathSnapped = "your_file_0.excalidraw"

var snapCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snap the diagram to a grid",
	Long: `This command is used to tidy up a diagram by snapping it's objects to a grid.

  This is an experimental feature.

Example:
  exconv snap -i your_file.excalidraw
`,
	Run: func(cmd *cobra.Command, args []string) {
		importPath, _ := cmd.Flags().GetString("input")
		exportPath, _ := cmd.Flags().GetString("output")
		gridSize, _ := cmd.Flags().GetString("grid-size")

		if len(importPath) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Input file path not provided.\n\n")
			cmd.Help()
			os.Exit(1)
		}

		gridSizeFloat, err := strconv.ParseFloat(gridSize, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse grid size: %s\n", err)
			os.Exit(1)
		}

		err = snap.SnapExcalidrawDiagramToGridAndSaveToFile(importPath, exportPath, gridSizeFloat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to snap Excalidraw diagram to grid: %s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(snapCmd)

	snapCmd.PersistentFlags().StringP("input", "i", "", "input file path")
	snapCmd.PersistentFlags().StringP("output", "o", defaultOutputPathSnapped, "output file path")
	snapCmd.PersistentFlags().StringP("grid-size", "g", "20", "grid size")
}
