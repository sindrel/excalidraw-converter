package internal

import (
	"io"
	"math"
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
