package internal

import (
	"errors"
	"runtime"
)

func DetectArch() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64", nil
	case "arm64":
		return "arm64", nil
	case "386":
		return "i386", nil
	default:
		return "", errors.New("unsupported architecture")
	}
}
