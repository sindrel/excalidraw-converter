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

	ignoredArrowAnchors := make(map[string]struct{})

	// First pass: snap non-container elements and record offsets for their bound elements
	for i, el := range input.Elements {
		if el.ContainerId != "" {
			// Skip container children in this pass
			continue
		}

		if len(el.GroupIds) > 0 {
			// Skip elements that are part of a group
			ignoredArrowAnchors[el.ID] = struct{}{}
			continue
		}

		// Snap element size to grid
		newWidth, newHeight := getSnappedElementSize(el.Width, el.Height, gridSize)
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
		newX, newY := getSnappedGridPosition(el.X, el.Y, gridSize)
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

		if len(el.GroupIds) > 0 {
			// Skip elements that are part of a group
			ignoredArrowAnchors[el.ID] = struct{}{}
			continue
		}

		// Apply size and position offsets if present
		output.Elements[i].Width = el.Width + sizeOffsets[el.ID].Width
		output.Elements[i].Height = el.Height + sizeOffsets[el.ID].Height
		output.Elements[i].X = el.X + positionOffsets[el.ID].X
		output.Elements[i].Y = el.Y + positionOffsets[el.ID].Y
	}

	// Third pass: adjust arrows to point to correct bound elements after snapping
	// Build a map from element ID to its new position and size
	elemPosSize := make(map[string]struct{ X, Y, W, H float64 })
	for _, el := range output.Elements {
		elemPosSize[el.ID] = struct{ X, Y, W, H float64 }{el.X, el.Y, el.Width, el.Height}
	}

	for i, el := range output.Elements {
		if el.Type != "arrow" {
			continue
		}

		// Skip arrows whose start or end bindings are in ignoredArrowAnchors
		if (el.StartBinding.ElementID != "" && func() bool { _, ok := ignoredArrowAnchors[el.StartBinding.ElementID]; return ok }()) ||
			(el.EndBinding.ElementID != "" && func() bool { _, ok := ignoredArrowAnchors[el.EndBinding.ElementID]; return ok }()) {
			continue
		}

		// Adjust start point if StartBinding is present
		if el.StartBinding.ElementID != "" {
			if bound, ok := elemPosSize[el.StartBinding.ElementID]; ok {
				// Use original start point and bound rect to determine side
				origStartX := el.X
				origStartY := el.Y
				side := getClosestSide(origStartX, origStartY, bound.X, bound.Y, bound.W, bound.H)
				newStartX, newStartY := getSideCoords(side, bound.X, bound.Y, bound.W, bound.H)
				output.Elements[i].X = newStartX
				output.Elements[i].Y = newStartY
			}
		}
		// Adjust end point if EndBinding is present
		if el.EndBinding.ElementID != "" && len(el.Points) > 0 {
			if bound, ok := elemPosSize[el.EndBinding.ElementID]; ok {
				// Use original end point and bound rect to determine side
				origEndX := el.X + el.Points[len(el.Points)-1][0]
				origEndY := el.Y + el.Points[len(el.Points)-1][1]
				side := getClosestSide(origEndX, origEndY, bound.X, bound.Y, bound.W, bound.H)
				newEndX, newEndY := getSideCoords(side, bound.X, bound.Y, bound.W, bound.H)
				startX := output.Elements[i].X
				startY := output.Elements[i].Y
				output.Elements[i].Points[len(el.Points)-1][0] = newEndX - startX
				output.Elements[i].Points[len(el.Points)-1][1] = newEndY - startY
			}
		}
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		return "", errors.New("Error occurred during JSON marshaling + " + err.Error())
	}

	return string(outputJson), nil
}

// Helper to determine which side of a rectangle a point is closest to
func getClosestSide(px, py, rectX, rectY, rectW, rectH float64) string {
	dxLeft := math.Abs(px - rectX)
	dxRight := math.Abs(px - (rectX + rectW))
	dyTop := math.Abs(py - rectY)
	dyBottom := math.Abs(py - (rectY + rectH))

	minDist := dxLeft
	side := "left"
	if dxRight < minDist {
		minDist = dxRight
		side = "right"
	}
	if dyTop < minDist {
		minDist = dyTop
		side = "top"
	}
	if dyBottom < minDist {
		// minDist = dyBottom // not needed
		side = "bottom"
	}
	return side
}

// Helper to get the coordinates for a given side of a rectangle
func getSideCoords(side string, rectX, rectY, rectW, rectH float64) (float64, float64) {
	switch side {
	case "left":
		return rectX, rectY + rectH/2
	case "right":
		return rectX + rectW, rectY + rectH/2
	case "top":
		return rectX + rectW/2, rectY
	case "bottom":
		return rectX + rectW/2, rectY + rectH
	default:
		return rectX + rectW/2, rectY + rectH/2 // fallback to center
	}
}

func getSnappedElementSize(width, height, gridSize float64) (float64, float64) {
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

func getSnappedGridPosition(xPos, yPos, gridSize float64) (float64, float64) {
	var snappedXPos = math.Ceil(xPos / gridSize)
	var snappedYPos = math.Ceil(yPos / gridSize)

	snappedXPos *= gridSize
	snappedYPos *= gridSize

	return snappedXPos, snappedYPos
}
