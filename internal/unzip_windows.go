//go:build windows

package internal

import (
	"archive/zip"
	"io"
	"strings"
)

func (u Unzip) decompress(dst io.Writer, src string) error {
	var isDecompressed = false

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		if !strings.Contains(f.Name, "client") || f.FileInfo().IsDir() {
			continue
		}

		_, err = io.Copy(dst, rc)
		if err != nil {
			return err
		}

		isDecompressed = true

		break // we just want to decompress one file which is named client, the one being executable.
	}

	if !isDecompressed {
		return ErrExeNotDecompressed
	}

	return nil
}
