package internal

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	WindowsBinaryExtension = "exe"
)

func DetectOS() (string, error) {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "Darwin", nil
	case "linux":
		return "Linux", nil
	case "windows":
		return "Windows", nil
	default:
		return "", errors.New("unknown OS")
	}
}

func FilenameWithExtension(filename string) (string, error) {
	os, err := DetectOS()

	if err != nil {
		return "", err
	}

	if os == "Windows" {
		return fmt.Sprintf("%s.%s", filename, WindowsBinaryExtension), nil
	}

	return filename, nil // unix binaries have not any extension
}
