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
	isLatest, _, err := checkIfLatestVersion(githubRepoUser, githubRepoName, version)

	assert.Nil(t, err)
	assert.Equal(t, isLatest, false)
}

func TestSanitizeElementText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & World", "Hello &amp; World"},
		{"<tag>", "&lt;tag&gt;"},
		{"Line1\nLine2", "Line1<br>Line2"},
		{"<b>&</b>", "&lt;b&gt;&amp;&lt;/b&gt;"},
		{"NoSpecialChars", "NoSpecialChars"},
		{"&<>", "&amp;&lt;&gt;"},
		{"Line1\nLine2\nLine3", "Line1<br>Line2<br>Line3"},
	}

	for _, tt := range tests {
		result := SanitizeElementText(tt.input)
		assert.Equal(t, tt.expected, result, "input: %q", tt.input)
	}
}
