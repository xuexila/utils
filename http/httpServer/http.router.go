package httpServer

import (
	"bytes"
	"fmt"
	"github.com/xuexila/utils/close/httpClose"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// 显示 favicon
func (ro Router) favicon(w http.ResponseWriter) {
	w.WriteHeader(200)
	rd := bytes.NewReader(favicon[:])
	_, _ = io.Copy(w, rd)
}

// Index 默认页面
func (ro Router) Index(w http.ResponseWriter, r *http.Request) {
	files, path, status := SetRequestDefaultPage(ro.Default, ro.Root+r.URL.String())
	defer func() {
		if files != nil {
			for _, item := range files {
				if item != nil {
					_ = item.Close()
				}
			}
		}
		httpClose.CloseReq(r)
	}()
	var rmime string
	if len(path) < 1 {
		rmime = "text/html; charset=utf-8"
	} else {
		if len(filepath.Ext(path[0])) > 0 {
			rmime = MimeMap[strings.ToLower(filepath.Ext(path[0])[1:])]
		}
		if rmime == "" {
			rmime = "text/html; charset=utf-8"
		}
	}
	w.Header().Set("Content-Type", rmime)
	if !status {
		if r.URL.Path == "/favicon.ico" {
			ro.favicon(w)
			return
		}
		NotFound(w, "404 not found")
		return
	}
	if r.Method == "POST" {
		MethodNotAllow(w)
		return
	}

	w.WriteHeader(200)
	if ro.HttpCache {
		w.Header().Set("cache-control", "max-age="+ro.HttpCacheMaxAge)
		if len(files) == 1 {
			fileInfo, _ := files[0].Stat()
			w.Header().Set("last-modified", fileInfo.ModTime().Format(time.RFC822))
		}
	}
	for _, file := range files {
		_, _ = io.Copy(w, file)
		_, _ = fmt.Fprintln(w)
	}
}

// 中间件
func (ro *Router) Middleware(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request, ro *Router)) {
	f(w, r, ro)
}
