package common

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type HttpServer struct {
	ListenAddr     string        `ini:"listen_addr"`
	Auth           string        `ini:"auth"`
	Allowip        []string      `ini:"allowip,omitempty"`
	Denyip         []string      `ini:"denyip,omitempty"`
	Ssl            bool          `ini:"ssl"`
	Ca             string        `ini:"ca"`
	Crt            string        `ini:"crt"`
	Key            string        `ini:"key"`
	SocketTimeout  time.Duration `ini:"socket_timeout"` // socket 心跳超时时间
	Hotupdate      bool          `ini:"hotupdate"`      // 是否启动热加载
	Route          map[string]func(w http.ResponseWriter, r *http.Request)
	RouteSocket    map[string]func(ws *websocket.Conn)
	CommonCallback func(w http.ResponseWriter, r *http.Request) bool
}

type Router struct {
	Default         string `ini:"default"`
	Root            string `ini:"root"`
	HttpCache       bool   `ini:"http_cache"`
	HttpCacheMaxAge string `ini:"http_cache_max_age"`
	SessionId       string `ini:"session_id"`
	CookiePath      string `ini:"cookie_path"`
	CookieDomain    string `ini:"cookie_domain"`
	CookieSecure    bool   `ini:"cookie_secure"`
	CookieHttpOnly  bool   `ini:"cookie_http_only"`
}
