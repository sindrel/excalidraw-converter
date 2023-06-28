package cmd

import (
	conv "diagram-converter/internal/conversion"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
	gliffyCmd.PersistentFlags().StringP("output", "o", "output.gliffy", "output file path")
}
