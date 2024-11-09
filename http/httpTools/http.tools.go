package httpTools

import (
	"net/http"
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
