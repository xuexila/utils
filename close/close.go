package close

import (
	"io"
)

// Close 通用关闭函数，不用操心有没有报错的情况
func Close(ide io.Closer) {
	if ide != nil {
		_ = ide.Close()
	}
}
