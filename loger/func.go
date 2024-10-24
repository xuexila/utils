package loger

import "os"

// 判断目录是否存在，否则创建目录
func mkdir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

