package internal

import (
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func NormalizeRotation(angle float64) float64 {
	return (angle * 180) / math.Pi
}

func checkIfLatestVersion(user, repo, version string) (bool, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", user, repo)

	resp, err := http.Get(url)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, "", errors.New("could not fetch releases")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	var release datastr.GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return false, "", err
	}

	return release.TagName == version, release.TagName, nil
}

func PrintVersionCheck(user, repo, version string) error {
	isLatest, latest, err := checkIfLatestVersion(user, repo, version)
	if err != nil {
		return err
	}

	if !isLatest {
		fmt.Printf("\nA newer version is available (%s). Go to 'https://github.com/%s/%s' for instructions on how to install the latest version.\n", latest, user, repo)
	} else {
		fmt.Printf("\nYou are using the latest version.\n")
	}

	return nil
}
