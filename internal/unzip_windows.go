//go:build windows

package internal

import "io"

type Unzip struct{}

func NewUnzip() Unzip {
	return Unzip{}
}

func (Unzip) Decompress(dst io.Writer, src io.Reader) error {
	return nil
}
