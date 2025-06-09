package snapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSnappedGridPosition(t *testing.T) {
	xPos, yPos := getSnappedGridPosition(817, 523, 20)

	assert.Equal(t, float64(820), xPos)
	assert.Equal(t, float64(540), yPos)
}

func TestGetSnappedElementSize(t *testing.T) {
	width, height := getSnappedElementSize(817, 523, 20)

	assert.Equal(t, float64(820), width)
	assert.Equal(t, float64(520), height)
}
