package multipart_close

import "mime/multipart"

func CloseMultipartWriter(w *multipart.Writer) {
	if w != nil {
		_ = w.Close()
	}
}

func CloseMultipartFile(f multipart.File) {
	if f != nil {
		_ = f.Close()
	}
}
