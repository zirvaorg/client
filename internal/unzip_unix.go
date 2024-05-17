//go:build aix || dragonfly || freebsd || (js && wasm) || wasip1 || linux || netbsd || openbsd || solaris || darwin

package internal

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"strings"
)

type Unzip struct{}

func NewUnzip() Unzip {
	return Unzip{}
}

func (Unzip) Decompress(dst io.Writer, src io.Reader) error {
	zipStream, err := gzip.NewReader(src)
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

		if !strings.Contains(header.Name, "client") {
			continue
		}

		if header.Typeflag == tar.TypeReg {
			_, err := io.Copy(dst, tarReader)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
