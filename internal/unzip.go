package internal

import "io"

type Decompressor interface {
	Decompress(dst io.Writer, src io.Reader) error
}
