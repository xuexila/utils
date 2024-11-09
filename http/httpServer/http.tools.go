package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/helays/utils/close/osClose"
	"github.com/helays/utils/http/mime"
	"github.com/helays/utils/tools"
	"github.com/helays/utils/ulogs"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Play 公共函数文件
func Play(path string, w http.ResponseWriter, r *http.Request, args ...any) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer osClose.CloseFile(f)
	if err != nil {
		ulogs.Error("文件不存在", path)
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
	m := mime.MimeMap[fType]
	if m == "" {
		m = "text/html;charset=utf-8"
	}
	w.Header().Set("Content-Type", m)
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

// Forbidden 设置系统返回403
func Forbidden(w http.ResponseWriter, msg ...string) {
	w.WriteHeader(http.StatusForbidden)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>403 Forbidden!</title></head><body><h3 style="text-align:center">` + strings.Join(msg, " ") + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// NotFound 设置返回 404
func NotFound(w http.ResponseWriter, msg ...string) {
	w.WriteHeader(http.StatusNotFound)
	_html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>404 Not Found!</title></head><body><h3 style="text-align:center">` + strings.Join(msg, " ") + `</h3></body></html>`
	_, _ = fmt.Fprintln(w, _html)
	return
}

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`<h1 style="text-align:center;">405 Error!</h1>
<p style="text-align:center;">` + http.StatusText(http.StatusMethodNotAllowed) + `</p>`))
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`<h1 style="text-align:center;">500 Error!</h1>
<p style="text-align:center;">` + http.StatusText(http.StatusInternalServerError) + `</p>`))
}

func ReqError(r *http.Request, i ...any) {
	log.SetPrefix("")
	log.SetOutput(os.Stderr)
	var msg = []any{Getip(r), r.URL.String()}
	log.Println(append(msg, i...)...)
}

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
	ulogs.Checkerr(json.NewEncoder(w).Encode(resp), "SetReturnCode")
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
	ulogs.Checkerr(json.NewEncoder(w).Encode(map[string]any{
		"code": code,
		"msg":  msg[0],
	}), "SetReturn")
}

func SetReturnData(w http.ResponseWriter, code int, msg any, data ...any) {
	RespJson(w)
	r := resp{
		"code": code,
		"msg":  msg,
	}
	if len(data) == 1 {
		r["data"] = data[0]
	} else if len(data) > 1 {
		r["data"] = data
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(r), "SetReturnData")
}

// SetReturnFile 直接讲文件反馈给前端
func SetReturnFile(w http.ResponseWriter, r *http.Request, file string) {
	f, err := os.Open(file)
	defer osClose.CloseFile(f)
	if err != nil {
		SetReturnError(w, r, err, http.StatusForbidden, "模板下载失败")
	}
	// 设置响应头
	mimeType, _ := mime.GetFilePathMimeType(file)
	w.Header().Set("Content-Type", mimeType)
	// 对文件名进行URL转义，以支持中文等非ASCII字符
	fileName := filepath.Base(file)
	encodedFileName := url.QueryEscape(fileName)
	// 设置Content-Disposition头部，使用RFC5987兼容的方式指定文件名
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s", encodedFileName, encodedFileName)
	w.Header().Set("Content-Disposition", contentDisposition)
	_, _ = io.Copy(w, f)
}

func SetReturnError(w http.ResponseWriter, r *http.Request, err error, code int, msg ...any) {
	ReqError(r, append([]any{err}, msg...)...)
	if len(msg) < 1 {
		msg = []any{err.Error()}
	} else {
		msg = append(msg, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	if code == 0 || code == 200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(code)
	}
	ulogs.Checkerr(json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"msg":  tools.AnySlice2Str(msg),
	}), "SetReturnError")
}

// CheckReqPost 检查请求是否post
func CheckReqPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		Forbidden(w, "Forbidden")
		return false
	}
	return true
}

// CheckSessionActive
// 检测session 有效期，超时就删除
func CheckSessionActive() {
	for range time.NewTicker(1 * time.Hour).C {
		LoginSessionMap.Range(func(key, value interface{}) bool {
			sid := key.(string)
			val := value.(LoginInfo)
			if val.HoldTime < 1 {
				val.HoldTime = 604800 // 默认保留7填
			}
			if int(time.Since(val.ActiveTime).Seconds()) > val.HoldTime {
				LoginSessionMap.Delete(sid)
			}
			return true
		})
	}
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
