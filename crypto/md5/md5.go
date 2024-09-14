package md5

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// Md5 md5 函数
func Md5(s []byte) string {
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

// Md5string 给字符串Md5
func Md5string(s string) string {
	return Md5([]byte(s))
}

// Md5file 计算文件的Md5
func Md5file(path string) string {
	byt, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return Md5(byt)
}
