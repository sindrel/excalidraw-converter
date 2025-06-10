package cmd

import (
	conv "diagram-converter/internal/conversion"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertExcalidrawToMermaidFile(t *testing.T) {
	inputDataPath := "../test/data/test_input_mermaid.excalidraw"
	outputDataPath := "../test/data/test_output.mermaid.tmp"
	expectedDataPath := "../test/data/test_output.mermaid"

	err := conv.ConvertExcalidrawDiagramToMermaidAndSaveToFile(inputDataPath, outputDataPath)
	assert.Nil(t, err)

	outputData, err := os.ReadFile(outputDataPath)
	assert.Nil(t, err)

	expectedData, err := os.ReadFile(expectedDataPath)
	assert.Nil(t, err)

	output := string(outputData)
	expected := string(expectedData)

	assert.Equal(t, output, expected)
}
