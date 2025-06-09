package cmd

import (
	snap "diagram-converter/internal/snapping"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var snapCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snap a diagram to grid",
	Long: `This command is used to tidy up an Excalidraw diagram by snapping it's objects to a grid.

	Resizes and aligns diagram objects to a grid. This can be useful to 
	quickly clean up sketches that are out of alignment. Objects will be
	resized to fit the grid and placed along the lines of the grid.
	
	It can also be used in-line when using the 'gliffy' command.

Example:
  exconv snap -i your_file.excalidraw
`,
	Run: func(cmd *cobra.Command, args []string) {
		importPath, _ := cmd.Flags().GetString("input")
		exportPath, _ := cmd.Flags().GetString("output")
		gridSize, _ := cmd.Flags().GetString("grid-size")
		forced, _ := cmd.Flags().GetBool("yes")

		if len(importPath) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Input file path not provided.\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if len(exportPath) == 0 {
			exportPath = importPath

			if !forced {
				fmt.Printf("No output file path supplied. The input file '%s' will be overwritten. Continue? (y/N): ", importPath)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Operation cancelled.")
					os.Exit(0)
				}
			}
		}

		gridSizeInt, err := strconv.ParseInt(gridSize, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse grid size: %s\n", err)
			os.Exit(1)
		}
		gridSizeFloat := float64(gridSizeInt)

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
	snapCmd.PersistentFlags().StringP("output", "o", "", "output file path (default: overwrite input file)")
	snapCmd.PersistentFlags().StringP("grid-size", "g", "20", "grid size")
	snapCmd.PersistentFlags().BoolP("yes", "y", false, "auto-approve overwriting the input file")

}
