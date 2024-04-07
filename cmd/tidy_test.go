package cmd

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	snap "diagram-converter/internal/snapping"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	sizeOffsets := make(map[string]datastr.ElementSizeOffset)
	positionOffsets := make(map[string]datastr.ElementPositionOffset)

	for i := range input.Elements {
		if input.Elements[i].ContainerId != "" {
			continue
		}

		newWidth, newHeight := snap.GetSnappedElementSize(input.Elements[i].Width, input.Elements[i].Height, gridSize)
		sizeDiffWidth := newWidth - input.Elements[i].Width
		sizeDiffHeight := newHeight - input.Elements[i].Height

		for _, boundElement := range input.Elements[i].BoundElements {
			sizeOffsets[boundElement.ID] = datastr.ElementSizeOffset{
				Width:  sizeDiffWidth,
				Height: sizeDiffHeight,
			}
		}

		newX, newY := snap.GetSnappedGridPosition(input.Elements[i].X, input.Elements[i].Y, gridSize)
		positionDiffX := newX - input.Elements[i].X
		positionDiffY := newY - input.Elements[i].Y

		for _, boundElement := range input.Elements[i].BoundElements {
			positionOffsets[boundElement.ID] = datastr.ElementPositionOffset{
				X: positionDiffX,
				Y: positionDiffY,
			}
		}

		// for element_k, element_v := range input.Elements[i] {
		// 	fmt.Println(element_k, element_v)
		// }

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
