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
func SnapExcalidrawDiagramToGridAndSaveToFile(inputPath string, outputPath string, gridSize float64) error {
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

func SnapExcalidrawDiagramToGrid(data string, gridSize float64) (string, error) {
	fmt.Printf("Aligning diagram elements to grid...\n")
	fmt.Printf("Grid size is: %f\n", gridSize)

	var input datastr.ExcalidrawScene
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		return "", errors.New("Unable to parse input: " + err.Error())
	}

	output := input
	output.AppState.GridSize = int64(gridSize)

	// These maps track how much each element's size and position should be offset
	sizeOffsets := make(map[string]datastr.ElementSizeOffset)
	positionOffsets := make(map[string]datastr.ElementPositionOffset)

	// First pass: snap non-container elements and record offsets for their bound elements
	for i, el := range input.Elements {
		if el.ContainerId != "" {
			// Skip container children in this pass
			continue
		}

		// Snap element size to grid
		newWidth, newHeight := GetSnappedElementSize(el.Width, el.Height, gridSize)
		sizeDiffWidth := newWidth - el.Width
		sizeDiffHeight := newHeight - el.Height

		// Record size offset for all bound elements (usually container children)
		for _, boundElement := range el.BoundElements {
			sizeOffsets[boundElement.ID] = datastr.ElementSizeOffset{
				Width:  sizeDiffWidth,
				Height: sizeDiffHeight,
			}
		}

		// Snap element position to grid
		newX, newY := GetSnappedGridPosition(el.X, el.Y, gridSize)
		positionDiffX := newX - el.X
		positionDiffY := newY - el.Y

		// Record position offset for all bound elements
		for _, boundElement := range el.BoundElements {
			positionOffsets[boundElement.ID] = datastr.ElementPositionOffset{
				X: positionDiffX,
				Y: positionDiffY,
			}
		}

		output.Elements[i].Width = newWidth
		output.Elements[i].Height = newHeight
		output.Elements[i].X = newX
		output.Elements[i].Y = newY
	}

	// Second pass: apply recorded offsets to container children
	for i, el := range input.Elements {
		if el.ContainerId == "" {
			// Only process container children in this pass
			continue
		}

		// Apply size and position offsets if present
		output.Elements[i].Width = el.Width + sizeOffsets[el.ID].Width
		output.Elements[i].Height = el.Height + sizeOffsets[el.ID].Height
		output.Elements[i].X = el.X + positionOffsets[el.ID].X
		output.Elements[i].Y = el.Y + positionOffsets[el.ID].Y
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
