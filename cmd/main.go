package main

import (
	conv "diagram-converter/internal/conversion"
	datastr "diagram-converter/internal/datastructures"
	"fmt"
	"os"
)

var version string
var commit string

var graphics datastr.GraphicTypes

func printHelp() {
	fmt.Printf("Usage:\n%s gliffy <input.excalidraw> <output.gliffy>\n\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Unexpected number of arguments\n")
		printHelp()
		os.Exit(1)
	}

	var command = os.Args[1]

	if command == "help" {
		printHelp()
		os.Exit(0)
	}

	if command == "version" {
		fmt.Printf("v%s (%s)\n", version, string(commit[0:7]))
		os.Exit(0)
	}

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Unexpected number of arguments\n")
		printHelp()
		os.Exit(1)
	}

	var importPath = os.Args[2]
	var exportPath = os.Args[3]

	if command != "gliffy" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}

	err := conv.ConvertExcalidrawToGliffy(importPath, exportPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to convert Excalidraw diagram to Gliffy diagram: %s\n", err)
		os.Exit(1)
	}
}
