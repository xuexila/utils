package tools

import "io"

func MustReadReader(r io.Reader) []byte {
	b, _ := io.ReadAll(r)
	return b
}

func MustStringReader(r io.Reader) string {
	b, _ := io.ReadAll(r)
	return string(b)
}
