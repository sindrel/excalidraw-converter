package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	version        = "v0.0.1"
	gitHubRepoUser = "sindrel"
	gitHubRepoName = "excalidraw-converter"
)

func TestCheckIfLatestVersion(t *testing.T) {
	isLatest, _, err := CheckIfLatestVersion(gitHubRepoUser, gitHubRepoName, version)

	assert.Nil(t, err)
	assert.Equal(t, isLatest, false)
}
