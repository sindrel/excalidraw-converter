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

var defaultOutputPath = "your_file.gliffy"

var gliffyCmd = &cobra.Command{
	Use:   "gliffy",
	Short: "Convert an Excalidraw diagram to Gliffy format",
	Long: `This command is used to convert an Excalidraw diagram to the Gliffy format.

  When an output file path is not provided, it will be determined 
automatically based on the filename of the input file. I.e. the 
input file path './subfolder/your_file.excalidraw' will produce 
the default output file path './your_file.gliffy'.

Example:
  exconv gliffy -i your_file.excalidraw
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
			exportPath = strings.TrimSuffix(path.Base(importPath), filepath.Ext(importPath)) + ".gliffy"
		}

		err := conv.ConvertExcalidrawToGliffyFile(importPath, exportPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to convert Excalidraw diagram to Gliffy diagram: %s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(gliffyCmd)

	gliffyCmd.PersistentFlags().StringP("input", "i", "", "input file path")
	gliffyCmd.PersistentFlags().StringP("output", "o", defaultOutputPath, "output file path")
}
