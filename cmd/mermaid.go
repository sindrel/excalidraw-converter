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

var defaultOutputPathMermaid = "your_file.mermaid"

var mermaidCmd = &cobra.Command{
	Use:   "mermaid",
	Short: "Convert an Excalidraw diagram to Mermaid format",
	Long: `This command is used to convert an Excalidraw diagram to the Mermaid format.

  When an output file path is not provided, it will be determined
automatically based on the filename of the input file. I.e. the
input file path './subfolder/your_file.excalidraw' will produce
the default output file path './your_file.mermaid'.

The diagram will be interpreted as a Mermaid flowchart. 
Only elements that are connected by arrows, or contained inside 
a connected element, will be included in the output.

Example:
  exconv mermaid -i your_file.excalidraw
`,
	Run: func(cmd *cobra.Command, args []string) {
		importPath, _ := cmd.Flags().GetString("input")
		exportPath, _ := cmd.Flags().GetString("output")
		printToStdOut, _ := cmd.Flags().GetBool("print-to-stdout")
		flowDirection, _ := cmd.Flags().GetString("direction")

		if len(importPath) == 0 {
			fmt.Fprintf(os.Stderr, "Error: Input file path not provided.\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if printToStdOut {
			output, err := conv.ConvertExcalidrawDiagramToMermaidAndOutputAsString(importPath, exportPath, flowDirection)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to convert Excalidraw diagram to Mermaid diagram: %s\n", err)
				os.Exit(1)
			}
			fmt.Print("---\n", output)
			return
		}

		if exportPath == defaultOutputPathMermaid {
			exportPath = strings.TrimSuffix(path.Base(importPath), filepath.Ext(importPath)) + ".mermaid"
		}

		err := conv.ConvertExcalidrawDiagramToMermaidAndSaveToFile(importPath, exportPath, flowDirection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to convert Excalidraw diagram to Mermaid diagram: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mermaidCmd)

	mermaidCmd.PersistentFlags().StringP("direction", "d", "default", "flow direction 'default', 'top-down', 'left-right', 'right-left' or 'bottom-top'")
	mermaidCmd.PersistentFlags().StringP("input", "i", "", "input file path")
	mermaidCmd.PersistentFlags().StringP("output", "o", defaultOutputPathMermaid, "output file path")
	mermaidCmd.PersistentFlags().BoolP("print-to-stdout", "p", false, "print output to stdout instead of a file")
}
