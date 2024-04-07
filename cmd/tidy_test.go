package cmd

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetSnappedGridPosition(t *testing.T) {
	xPos, yPos := GetSnappedGridPosition(817, 523, 20)

	assert.Equal(t, float64(820), xPos)
	assert.Equal(t, float64(540), yPos)
}

func TestGetSnappedElementSize(t *testing.T) {
	width, height := GetSnappedElementSize(817, 523, 20)

	assert.Equal(t, float64(820), width)
	assert.Equal(t, float64(520), height)
}

type elementSizeOffset struct {
	Width  float64
	Height float64
}

type elementPositionOffset struct {
	X float64
	Y float64
}

func TestSnapExcalidrawDiagram(t *testing.T) {
	inputDataPath := "../test/data/test_input.excalidraw"
	outputDataPath := "../test/data/test_output_snapped.excalidraw"
	// expectedDataPath := "../test/data/test_output.gliffy"

	inputData, err := os.ReadFile(inputDataPath)
	assert.Nil(t, err)

	var input datastr.ExcalidrawScene
	err = json.Unmarshal(inputData, &input)
	assert.Nil(t, err)

	gridSize := float64(20)

	// var output datastr.ExcalidrawScene
	output := input

	sizeOffsets := make(map[string]elementSizeOffset)
	positionOffsets := make(map[string]elementPositionOffset)

	for i := range input.Elements {
		if input.Elements[i].ContainerId != "" {
			continue
		}

		newWidth, newHeight := GetSnappedElementSize(input.Elements[i].Width, input.Elements[i].Height, gridSize)
		sizeDiffWidth := newWidth - input.Elements[i].Width
		sizeDiffHeight := newHeight - input.Elements[i].Height

		for _, boundElement := range input.Elements[i].BoundElements {
			sizeOffsets[boundElement.ID] = elementSizeOffset{
				Width:  sizeDiffWidth,
				Height: sizeDiffHeight,
			}
		}

		newX, newY := GetSnappedGridPosition(input.Elements[i].X, input.Elements[i].Y, gridSize)
		positionDiffX := newX - input.Elements[i].X
		positionDiffY := newY - input.Elements[i].Y

		for _, boundElement := range input.Elements[i].BoundElements {
			positionOffsets[boundElement.ID] = elementPositionOffset{
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
	assert.Nil(t, err)

	fmt.Println(positionOffsets)
	fmt.Println(sizeOffsets)

	internal.WriteToFile(outputDataPath, string(outputJson))

	// assert.Equal(t, output, expected)
}
