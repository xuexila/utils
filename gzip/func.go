package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
)

func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.LittleEndian, data); err != nil {
		return nil, nil
	}

	r, err := gzip.NewReader(b)
	if err != nil {
		// fmt.Println("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return data, nil
	} else {
		defer func() {
			_ = r.Close()
		}()
		undatas, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return undatas, nil
	}
}
