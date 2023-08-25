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
	Long:  `Use this command to convert an Excalidraw diagram to the Gliffy format.`,
	Run: func(cmd *cobra.Command, args []string) {
		importPath, _ := cmd.Flags().GetString("input")
		exportPath, _ := cmd.Flags().GetString("output")

		if len(importPath) == 0 {
			fmt.Fprintf(os.Stderr, "Input file path not provided. (Use --help for details.)\n")
			os.Exit(1)
		}

		if strings.HasPrefix(exportPath, defaultOutputPath) {
			exportPath = strings.TrimSuffix(path.Base(importPath), filepath.Ext(importPath)) + ".gliffy"
		}

		err := conv.ConvertExcalidrawToGliffy(importPath, exportPath)
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
