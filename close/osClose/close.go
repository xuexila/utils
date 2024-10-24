package osClose

import "os"

func CloseFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}
