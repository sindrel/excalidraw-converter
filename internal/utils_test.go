package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	version        = "v0.0.1"
	githubRepoUser = "sindrel"
	githubRepoName = "excalidraw-converter"
)

func TestCheckIfLatestVersion(t *testing.T) {
	isLatest, _, err := CheckIfLatestVersion(githubRepoUser, githubRepoName, version)

	assert.Nil(t, err)
	assert.Equal(t, isLatest, false)
}
