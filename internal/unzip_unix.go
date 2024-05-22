//go:build aix || dragonfly || freebsd || (js && wasm) || wasip1 || linux || netbsd || openbsd || solaris || darwin

package internal

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

func (u Unzip) decompress(dst io.Writer, src string) error {
	var isDecompressed = false

	compressedFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(compressedFile *os.File) {
		_ = compressedFile.Close()
	}(compressedFile)

	zipStream, err := gzip.NewReader(compressedFile)
	if err != nil {
		return err
	}
	defer zipStream.Close()

	tarReader := tar.NewReader(zipStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if !strings.Contains(header.Name, "client") || header.FileInfo().IsDir() {
			continue
		}

		if header.Typeflag == tar.TypeReg {
			_, err := io.Copy(dst, tarReader)
			if err != nil {
				return err
			}
			isDecompressed = true
			break
		}
	}

	if !isDecompressed {
		return ErrExeNotDecompressed
	}

	return nil
}
