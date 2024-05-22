package internal

import (
	"fmt"
	"runtime"
)

const (
	WindowsBinaryExtension = "exe"
)

func FilenameWithExtension(filename string) string {
	if os := runtime.GOOS; os == "windows" {
		return fmt.Sprintf("%s.%s", filename, WindowsBinaryExtension)
	}

	return filename // unix binaries have not any extension
}
