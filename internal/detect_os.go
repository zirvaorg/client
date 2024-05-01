package internal

import (
	"errors"
	"runtime"
)

func DetectOS() (string, error) {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "macOS", nil
	case "linux":
		return "Linux", nil
	case "windows":
		return "Windows", nil
	default:
		return "", errors.New("unknown OS")
	}
}
