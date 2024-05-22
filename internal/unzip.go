package internal

import (
	"errors"
	"io"
)

type Decompressor interface {
	Decompress(dst io.Writer, src string) error
}

type Unzip struct{}

var (
	ErrExeNotDecompressed = errors.New("executable file not decompressed")
)

func NewUnzip() Unzip {
	return Unzip{}
}

func (u Unzip) Decompress(dst io.Writer, src string) error {
	return u.decompress(dst, src)
}
