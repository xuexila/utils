package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func RespJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetReturnCode(w http.ResponseWriter, r *http.Request, code int, msg any, data ...any) {
	w.Header().Set("Content-Type", "application/json")
	if code == 0 || code == 200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(code)
		ReqError(r, code, msg)
	}

	if _, ok := msg.(error); ok {
		if len(data) > 0 && reflect.TypeOf(data[0]).String() == "bool" && !data[0].(bool) {
			msg = "系统处理失败"
		} else {
			msg = msg.(error).Error()
		}
	}
	resp := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if len(data) > 0 {
		resp["data"] = data
		if len(data) == 1 {
			resp["data"] = data[0]
		}
	}
	Checkerr(json.NewEncoder(w).Encode(resp), "SetReturnCode")
}

// SetReturn 设置 返回函数Play
func SetReturn(w http.ResponseWriter, code int, msg ...any) {
	w.Header().Set("Content-Type", "application/json")
	if len(msg) < 1 {
		if code == 0 {
			msg = []any{"成功"}
		} else {
			msg = []any{"失败"}
		}
	}
	Checkerr(json.NewEncoder(w).Encode(map[string]any{
		"code": code,
		"msg":  msg[0],
	}), "SetReturn")
}

func SetReturnData(w http.ResponseWriter, code int, msg any, data any) {
	w.Header().Set("Content-Type", "application/json")
	Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}), "SetReturnData")
}

func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	ReqError(r, append([]any{err}, msg...)...)
	w.Header().Set("Content-Type", "application/json")
	Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  msg,
	}), "SetReturnError")
}

// CheckReqPost 检查请求是否post
func CheckReqPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		Forbidden(w, "Forbidden")
		return false
	}
	return true
}

// Getip 获取客户端IP
func Getip(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("HTTP_CLIENT_IP"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// HttpServerMiddleware http统一验证 中间件
func HttpServerMiddleware(next http.Handler, callback func(w http.ResponseWriter, r *http.Request) bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r != nil && r.Body != nil {
				_ = r.Body.Close()
			}
		}()
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("server", "vs/1.0")
		if callback != nil && !callback(w, r) {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MethodNotAllowed 设置返回 405
func MethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
}

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusForbidden)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>403 Forbidden!</title></head><body><h3 style="text-align:center">` + msg + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// NotFound 设置返回 404
func NotFound(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusNotFound)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>404 Not Found!</title></head><body><h3 style="text-align:center">` + msg + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

func ReqError(r *http.Request, i ...any) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	var msg = []any{Getip(r), r.URL.String()}
	log.Println(append(msg, i...)...)
}

// Play 公共函数文件
func Play(path string, w http.ResponseWriter, r *http.Request, args ...any) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer CloseFile(f)
	if err != nil {
		Error("文件不存在", path)
		http.NotFound(w, r)
		return
	}
	ranges := int64(0)
	rangeEnd := int(0)
	rangeSwap := strings.TrimSpace(r.Header.Get("Range"))
	if rangeSwap != "" {
		rangeSwap = rangeSwap[6:]
		rangListSwap := strings.Split(rangeSwap, "-")
		if len(rangeSwap) >= 1 {
			if num, err := strconv.Atoi(rangListSwap[0]); err == nil {
				ranges = int64(num)
			}
			if len(rangeSwap) > 1 {
				if num, err := strconv.Atoi(rangListSwap[1]); err == nil {
					rangeEnd = int(num)
				}
			}
		}
	}

	var (
		fileSize int
		//tmpF     []byte
	)
	fType := strings.ToLower(filepath.Ext(path)[1:])
	fInfo, err := f.Stat()
	if err != nil {
		Forbidden(w, "403 Forbidden!")
		return
	}

	//if fType=="mp4" {
	//	tmpF = GetMP4Duration(f)
	//	fileSize= len(tmpF)
	//	if rangeSwap!="" && ranges>0 {
	//		tmpF=tmpF[ranges:]
	//	}
	//}else{
	//
	//	if rangeSwap!="" && ranges>0 {
	//		_, _ = f.Seek(ranges, 0)
	//	}
	//	fileSize=int(fInfo.Size())
	//}
	//GetMP4Duration(f)
	if rangeSwap != "" && ranges > 0 {
		_, _ = f.Seek(ranges, 0)
	}
	fileSize = int(fInfo.Size())

	totalSize := fileSize
	if rangeSwap != "" && rangeEnd > 0 {
		totalSize = rangeEnd
	}
	total := strconv.Itoa(fileSize)

	mime := MimeMap[fType]
	if mime == "" {
		mime = "text/html;charset=utf-8"
	}
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(totalSize-int(ranges)))
	w.Header().Set("Last-Modified", fInfo.ModTime().Format(time.RFC822))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Connection", "close")
	w.Header().Set("Etag", `W/"`+strconv.FormatInt(fInfo.ModTime().Unix(), 16)+`-`+strconv.FormatInt(fInfo.Size(), 16)+`"`)
	if len(args) > 0 {
		if args[0] == "downloader" {
			w.Header().Del("Accept-Ranges")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(path)))
		}
	}
	if rangeSwap != "" {
		w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(ranges))+"-"+strconv.Itoa(totalSize-1)+"/"+total)
		w.WriteHeader(206)
	} else {
		w.WriteHeader(200)
	}

	//if fType == "mp4" {
	//	_byt, _ := io.ReadAll(f)
	//	_, _ = w.Write(_byt)
	//	return
	//}
	_, _ = io.Copy(w, f)
}

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, `
<h1 style="text-align:center;">400 Error!</h1>
<p style="text-align:center;">`+http.StatusText(http.StatusMethodNotAllowed)+`</p>
`)
}
