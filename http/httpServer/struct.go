package httpServer

import (
	"golang.org/x/net/websocket"
	"net/http"
	"regexp"
	"time"
)

type HttpServer struct {
	ListenAddr     string                                                  `ini:"listen_addr" json:"listen_addr" yaml:"listen_addr"`
	Auth           string                                                  `ini:"auth" json:"auth" yaml:"auth"`
	Allowip        []string                                                `ini:"allowip,omitempty" json:"allowip" yaml:"allowip,omitempty"`
	Denyip         []string                                                `ini:"denyip,omitempty" json:"denyip" yaml:"denyip,omitempty"`
	ServerName     []string                                                `ini:"server_name,omitempty"` // 绑定域名
	Ssl            bool                                                    `ini:"ssl" json:"ssl" yaml:"ssl"`
	Ca             string                                                  `ini:"ca" json:"ca" yaml:"ca"`
	Crt            string                                                  `ini:"crt" json:"crt" yaml:"crt"`
	Key            string                                                  `ini:"key" json:"key" yaml:"key"`
	SocketTimeout  time.Duration                                           `ini:"socket_timeout" json:"socket_timeout" yaml:"socket_timeout"` // socket 心跳超时时间
	Hotupdate      bool                                                    `ini:"hotupdate" json:"hotupdate" yaml:"hotupdate"`                // 是否启动热加载
	EnableGzip     bool                                                    `ini:"enable_gzip" json:"enable_gzip" yaml:"enable_gzip"`          // 是否开启gzip
	Route          map[string]func(w http.ResponseWriter, r *http.Request) `yaml:"-" json:"-"`
	RouteSocket    map[string]func(ws *websocket.Conn)                     `yaml:"-" json:"-"`
	CommonCallback func(w http.ResponseWriter, r *http.Request) bool       `yaml:"-" json:"-"`
	serverNameMap  map[string]byte                                         // 绑定的域名
}

type Router struct {
	Default                string `ini:"default" json:"default" yaml:"default"`
	Root                   string `ini:"root" json:"root" yaml:"root"`
	HttpCache              bool   `ini:"http_cache" json:"http_cache" yaml:"http_cache"`
	HttpCacheMaxAge        string `ini:"http_cache_max_age" json:"http_cache_max_age" yaml:"http_cache_max_age"`
	UnauthorizedRespMethod int    `ini:"unauthorized_resp_method"` // 未登录响应方法 默认为 401，302 表示自动重定向到登录页面
	SessionId              string `ini:"session_id" json:"session_id" yaml:"session_id"`
	CookiePath             string `ini:"cookie_path" json:"cookie_path" yaml:"cookie_path"`
	CookieDomain           string `ini:"cookie_domain" json:"cookie_domain" yaml:"cookie_domain"`
	CookieSecure           bool   `ini:"cookie_secure" json:"cookie_secure" yaml:"cookie_secure"`
	CookieHttpOnly         bool   `ini:"cookie_http_only" json:"cookie_http_only" yaml:"cookie_http_only"`

	IsLogin                bool             // 是否登录
	LoginPath              string           // 登录页面
	HomePage               string           //首页
	UnLoginPath            map[string]bool  // 免授权页面
	UnLoginPathRegexp      []*regexp.Regexp // 免授权页面正则
	MustLoginPath          map[string]bool  //必须登录才能访问的页面
	MustLoginPathRegexp    []*regexp.Regexp // 必须登录才能访问的页面正则
	DisableLoginPath       map[string]bool  // 登录状态下不能访问的页面
	DisableLoginPathRegexp []*regexp.Regexp // 登录状态下不能访问的页面正则
	ManagePage             map[string]bool  // 管理员访问
	ManagePageRegexp       []*regexp.Regexp
}

// LoginInfo 登录信息
type LoginInfo struct {
	LoginTime     time.Time // 登录时间
	IsLogin       bool      // 是否登录
	UserId        int       // 用户ID
	User          string    // 用户名
	IsManage      bool      // 是否管理员
	DemoUser      bool      // 是否演示用户
	ActiveTime    time.Time // 最后活动时间
	HoldTime      int       // 会话保留时长
	RsaPrivateKey []byte    //ras 私钥
	RsaPublickKey []byte    // rsa 公钥
}
