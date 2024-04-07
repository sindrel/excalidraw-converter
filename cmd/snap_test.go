package cmd

import (
	datastr "diagram-converter/internal/datastructures"
	snap "diagram-converter/internal/snapping"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapExcalidrawDiagramToGridAndSaveToFile(t *testing.T) {
	inputDataPath := "../test/data/test_input.excalidraw"
	outputDataPath := "../test/data/test_output_snapped.excalidraw.tmp"
	expectedDataPath := "../test/data/test_output_snapped.excalidraw"

	err := snap.SnapExcalidrawDiagramToGridAndSaveToFile(inputDataPath, outputDataPath)
	assert.Nil(t, err)

	outputData, err := os.ReadFile(outputDataPath)
	assert.Nil(t, err)

	expectedData, err := os.ReadFile(expectedDataPath)
	assert.Nil(t, err)

	var output datastr.ExcalidrawScene
	err = json.Unmarshal(outputData, &output)
	assert.Nil(t, err)

	var expected datastr.ExcalidrawScene
	err = json.Unmarshal(expectedData, &expected)
	assert.Nil(t, err)

	assert.Equal(t, output, expected)
}
