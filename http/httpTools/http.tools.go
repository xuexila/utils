package httpTools

import (
	"fmt"
	"github.com/helays/utils/tools"
	"net/http"
	"net/url"
	"strings"
)

// QueryGetSlice 获取query参数，并分割
func QueryGetSlice(r *http.Request, key string, step string) []string {
	query := r.URL.Query()
	v := query.Get(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, step)
}

// QueryGet 获取query参数，如果值不存在就设置默认值
func QueryGet(query url.Values, k, dfValue string) string {
	v := query.Get(k)
	return tools.Ternary(v == "", dfValue, v)
}

// SetDisposition 文件下载时候，设置中文文件名
func SetDisposition(w http.ResponseWriter, filename string) {
	encodedFileName := url.QueryEscape(filename)
	// 设置Content-Disposition头部，使用RFC5987兼容的方式指定文件名
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName)
	w.Header().Set("Content-Disposition", contentDisposition)
}
