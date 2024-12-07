package sessionConfig

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
// User helay
// Date: 2024/12/7 22:34

const Interval = time.Hour // 默认检测频率
const CookieName = "vsclubId"
const ExpireTime = 24 * time.Hour // session默认24小时过期

var (
	ErrUnSupport = errors.New("不支持的session载体")
	ErrNotFound  = errors.New("session不存在")
)

func newSessionId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

// GetSessionId 获取sessionId
// 如果未获取到 sessionId ，则生成一个，并写到响应信息中
func GetSessionId(w http.ResponseWriter, r *http.Request, options *Options) (string, error) {
	switch options.Carrier {
	case "", "cookie":
		cookie, err := r.Cookie(options.CookieName)
		if err != nil {
			sid := newSessionId()
			SetSessionId(w, sid, options)
			return sid, nil
		}
		return cookie.Value, nil
	case "header":
		sid := r.Header.Get(options.CookieName)
		if sid == "" {
			sid = newSessionId()
			SetSessionId(w, sid, options)
		}
		return sid, nil
	}
	return "", ErrUnSupport

}

// SetSessionId 设置sessionId
func SetSessionId(w http.ResponseWriter, sid string, options *Options) {
	switch options.Carrier {
	case "cookie", "":
		http.SetCookie(w, &http.Cookie{
			Name:       options.CookieName,
			Value:      sid,
			Path:       options.Path,
			Domain:     options.Domain,
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     options.MaxAge,
			Secure:     options.Secure,
			HttpOnly:   options.HttpOnly,
			SameSite:   options.SameSite,
		})
	case "header":
		w.Header().Set("vsclub_"+options.CookieName, sid)
	}
}

func DeleteSessionId(w http.ResponseWriter, options *Options) {
	switch options.Carrier {
	case "cookie", "":
		http.SetCookie(w, &http.Cookie{
			Name:       options.CookieName,
			Value:      "",
			Path:       options.Path,
			Domain:     options.Domain,
			Expires:    time.Unix(0, 0),
			RawExpires: "",
			MaxAge:     -1,
			Secure:     options.Secure,
		})
	case "header":
		w.Header().Del("vsclub_" + options.CookieName)
	}
}

// GetSessionName 生成session key
// 这个key 是存储在 存储系统中的
func GetSessionName(sessionId, name string) string {
	return sessionId + "_" + name
}
