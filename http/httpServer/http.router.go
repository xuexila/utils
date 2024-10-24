package httpServer

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/helays/utils/close/httpClose"
	"github.com/helays/utils/crypto/md5"
	"github.com/helays/utils/http/mime"
	"github.com/helays/utils/ulogs"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
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
			rmime = mime.MimeMap[strings.ToLower(filepath.Ext(path[0])[1:])]
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

	if len(files) == 1 {
		fileInfo, _ := files[0].Stat()
		fileSize := int(fileInfo.Size())
		total := strconv.Itoa(fileSize)
		w.Header().Set("last-modified", fileInfo.ModTime().Format(time.RFC822))
		w.Header().Set("Accept-Ranges", "bytes")

		ranges := int64(0)
		rangeSwap := strings.TrimSpace(r.Header.Get("Range"))
		if rangeSwap != "" {
			rangeSwap = rangeSwap[6:]
			rangeListSwap := strings.Split(rangeSwap, "-")
			if len(rangeListSwap) == 2 {
				if num, err := strconv.Atoi(rangeListSwap[0]); err == nil {
					ranges = int64(num)
				}
			}
		}
		w.Header().Set("Content-Length", strconv.Itoa(fileSize-int(ranges)))
		_, _ = files[0].Seek(ranges, 0)
		w.Header().Set("Etag", `W/"`+strconv.FormatInt(fileInfo.ModTime().Unix(), 16)+`-`+strconv.FormatInt(fileInfo.Size(), 16)+`"`)

		if ranges > 0 {
			w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(ranges))+"-"+strconv.Itoa(fileSize-1)+"/"+total) // 允许 range
			w.WriteHeader(206)
		} else {
			w.WriteHeader(200)
		}
	} else {
		w.WriteHeader(200)
	}

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

// BeforeAction 所有应用前置操作
func (this *Router) BeforeAction(w http.ResponseWriter, r *http.Request) bool {
	if this.CookiePath == "" {
		this.CookiePath = "/"
	}
	this.IsLogin = true
	sessionId, err := GetSessionId(r, this.SessionId)
	if err != nil {
		// 未登录
		this.IsLogin = false
		this.SetCookie(w, this.SessionId, sessionId, "/")
		if _, ok := this.MustLoginPath[r.URL.Path]; ok {
			if this.UnauthorizedRespMethod == 401 {
				SetReturnCode(w, r, this.UnauthorizedRespMethod, "未登录，请先登录！！")
				return false
			}
			http.Redirect(w, r, this.LoginPath, 302)
			return false
		}
		return true
	}
	_loginInfo, ok := GetLoginInfo(sessionId)
	if !ok || !_loginInfo.IsLogin {
		this.IsLogin = false
		if _, ook := this.MustLoginPath[r.URL.Path]; ook {
			if this.UnauthorizedRespMethod == 401 {
				SetReturnCode(w, r, this.UnauthorizedRespMethod, "未登录，请先登录！！")
				return false
			}
			http.Redirect(w, r, this.LoginPath, 302)
			return false
		}
		return true
	}
	_loginInfo.ActiveTime = time.Now()
	LoginSessionMap.Store(sessionId, _loginInfo)
	if _, ok := this.DisableLoginPath[r.URL.Path]; ok {
		http.Redirect(w, r, this.HomePage, 302)
		return false
	}
	// 控制管理员访问的
	if _, ok := this.ManagePage[r.URL.Path]; ok && !_loginInfo.IsManage {
		SetReturnCode(w, r, http.StatusForbidden, "无权限访问")
		return false
	}
	return true
}

func (this Router) Captcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer
	sessionId, _ := GetSessionId(r, this.SessionId)
	captchaId := captcha.NewLen(4)
	this.SetCookie(w, "_c."+md5.Md5string(sessionId), captchaId, "/")
	if err := captcha.WriteImage(&content, captchaId, 106, 40); err != nil {
		InternalServerError(w)
		ulogs.Error(err, "captcha writeImage")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(content.Bytes()))
}

// 中间件
func (ro *Router) Middleware(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request, ro *Router)) {
	f(w, r, ro)
}
