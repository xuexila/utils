package utils_http

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type HttpServer struct {
	ListenAddr     string        `ini:"listen_addr" json:"listen_addr" yaml:"listen_addr"`
	Auth           string        `ini:"auth" json:"auth" yaml:"auth"`
	Allowip        []string      `ini:"allowip,omitempty" yaml:"allowip" json:"allowip"`
	Denyip         []string      `ini:"denyip,omitempty" json:"denyip" yaml:"denyip"`
	Ssl            bool          `ini:"ssl" json:"ssl" yaml:"ssl"`
	Ca             string        `ini:"ca" json:"ca" yaml:"ca"`
	Crt            string        `ini:"crt" json:"crt" yaml:"crt"`
	Key            string        `ini:"key" json:"key" yaml:"key"`
	SocketTimeout  time.Duration `ini:"socket_timeout" json:"socket_timeout" yaml:"socket_timeout"` // socket 心跳超时时间
	Hotupdate      bool          `ini:"hotupdate" json:"hotupdate" yaml:"hotupdate"`                // 是否启动热加载
	Route          map[string]func(w http.ResponseWriter, r *http.Request)
	RouteSocket    map[string]func(ws *websocket.Conn)
	CommonCallback func(w http.ResponseWriter, r *http.Request) bool
}

type Router struct {
	Default         string `ini:"default" json:"default" yaml:"default"`
	Root            string `ini:"root" json:"root" yaml:"root"`
	HttpCache       bool   `ini:"http_cache" json:"http_cache" yaml:"http_cache"`
	HttpCacheMaxAge string `ini:"http_cache_max_age" json:"http_cache_max_age" yaml:"http_cache_max_age"`
	SessionId       string `ini:"session_id" json:"session_id" yaml:"session_id"`
	CookiePath      string `ini:"cookie_path" json:"cookie_path" yaml:"cookie_path"`
	CookieDomain    string `ini:"cookie_domain" json:"cookie_domain" yaml:"cookie_domain"`
	CookieSecure    bool   `ini:"cookie_secure" json:"cookie_secure" yaml:"cookie_secure"`
	CookieHttpOnly  bool   `ini:"cookie_http_only" json:"cookie_http_only" yaml:"cookie_http_only"`
}
