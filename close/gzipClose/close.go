package gzipClose

import "compress/gzip"

func Close(rd *gzip.Reader) {
	if rd != nil {
		_ = rd.Close()
	}
}
