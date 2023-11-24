package cmd

import (
	conv "diagram-converter/internal/conversion"
	"encoding/json"
	"os"
	"testing"

	datastr "diagram-converter/internal/datastructures"

	"github.com/stretchr/testify/assert"
)

func TestConvertExcalidrawToGliffy(t *testing.T) {
	inputDataPath := "../test/data/test_input.excalidraw"
	outputDataPath := "../test/data/test_output.gliffy.tmp"
	expectedDataPath := "../test/data/test_output.gliffy"

	err := conv.ConvertExcalidrawToGliffyFile(inputDataPath, outputDataPath)
	assert.Nil(t, err)

	outputData, err := os.ReadFile(outputDataPath)
	assert.Nil(t, err)

	expectedData, err := os.ReadFile(expectedDataPath)
	assert.Nil(t, err)

	var output datastr.GliffyScene
	err = json.Unmarshal(outputData, &output)
	assert.Nil(t, err)

	var expected datastr.GliffyScene
	err = json.Unmarshal(expectedData, &expected)
	assert.Nil(t, err)

	expected.Metadata.LastSerialized = output.Metadata.LastSerialized

	assert.Equal(t, output, expected)
}
