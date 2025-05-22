package snapping

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
)

// TODO: This probably shouldn't be a separate function for snapped diagrams?
func SnapExcalidrawDiagramToGridAndSaveToFile(inputPath string, outputPath string, gridSize int64) error {
	fmt.Printf("Parsing input file: %s\n", inputPath)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File reading failed. %s\n", err)
		os.Exit(1)
	}

	output, err := SnapExcalidrawDiagramToGrid(string(data), gridSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File parsing failed. %s\n", err)
		os.Exit(1)
	}

	err = internal.WriteToFile(outputPath, string(output))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Saving diagram failed. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Snapped diagram saved to file: %s\n", outputPath)

	return nil
}

func SnapExcalidrawDiagramToGrid(data string, gridSize int64) (string, error) {
	fmt.Printf("Aligning diagram elements to grid...\n")
	fmt.Printf("Grid size is: %d\n", gridSize)

	var input datastr.ExcalidrawScene
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		return "", errors.New("Unable to parse input: " + err.Error())
	}

	gridSizeFloat := float64(gridSize)

	output := input

	sizeOffsets := make(map[string]datastr.ElementSizeOffset)
	positionOffsets := make(map[string]datastr.ElementPositionOffset)

	output.AppState.GridSize = gridSize

	for i := range input.Elements {
		if input.Elements[i].ContainerId != "" {
			continue
		}

		newWidth, newHeight := GetSnappedElementSize(input.Elements[i].Width, input.Elements[i].Height, gridSizeFloat)
		sizeDiffWidth := newWidth - input.Elements[i].Width
		sizeDiffHeight := newHeight - input.Elements[i].Height

		for _, boundElement := range input.Elements[i].BoundElements {
			sizeOffsets[boundElement.ID] = datastr.ElementSizeOffset{
				Width:  sizeDiffWidth,
				Height: sizeDiffHeight,
			}
		}

		newX, newY := GetSnappedGridPosition(input.Elements[i].X, input.Elements[i].Y, gridSizeFloat)
		positionDiffX := newX - input.Elements[i].X
		positionDiffY := newY - input.Elements[i].Y

		for _, boundElement := range input.Elements[i].BoundElements {
			positionOffsets[boundElement.ID] = datastr.ElementPositionOffset{
				X: positionDiffX,
				Y: positionDiffY,
			}
		}

		output.Elements[i].Width, output.Elements[i].Height = newWidth, newHeight
		output.Elements[i].X, output.Elements[i].Y = newX, newY
	}

	for i := range input.Elements {
		if input.Elements[i].ContainerId == "" {
			continue
		}

		output.Elements[i].Width, output.Elements[i].Height = input.Elements[i].Width+sizeOffsets[input.Elements[i].ID].Width, input.Elements[i].Height+sizeOffsets[input.Elements[i].ID].Height
		output.Elements[i].X, output.Elements[i].Y = input.Elements[i].X+positionOffsets[input.Elements[i].ID].X, input.Elements[i].Y+positionOffsets[input.Elements[i].ID].Y
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", errors.New("Error occurred during JSON marshaling + " + err.Error())
	}

	return string(outputJson), nil
}

func GetSnappedElementSize(width, height, gridSize float64) (float64, float64) {
	var adjustedWidth = math.Round(width / gridSize)
	var adjustedHeight = math.Round(height / gridSize)

	if adjustedWidth == 0 && width > 0 {
		adjustedWidth += 1
	}
	if adjustedHeight == 0 && height > 0 {
		adjustedHeight += 1
	}

	adjustedWidth *= gridSize
	adjustedHeight *= gridSize

	return adjustedWidth, adjustedHeight
}

func GetSnappedGridPosition(xPos, yPos, gridSize float64) (float64, float64) {
	var snappedXPos = math.Ceil(xPos / gridSize)
	var snappedYPos = math.Ceil(yPos / gridSize)

	snappedXPos *= gridSize
	snappedYPos *= gridSize

	return snappedXPos, snappedYPos
}
