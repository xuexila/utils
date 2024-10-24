package httpClient

import (
	"net/http"
	"net/url"
	"time"
)

type Curl struct {
	Curlcommon
	Type       string       // 请求类型 GET POST PUT。。。
	Client     *http.Client // 客户端
	acceptList []string     // 允许的content-type 集合
	Strict     bool
	jar        http.CookieJar
}

// Curlcommon
// Curl 公共配置
type Curlcommon struct {
	UserAgent      string `ini:"useragent"` // 系统默认
	Accept         string `ini:"accept"`
	AcceptLanguage string `ini:"language"`
	CacheControl   string `ini:"cache"`
	Connection     string `ini:"connection"`
	AcceptEncoding string `ini:"accept_encoding"` // 允许的传输方式
	Refer          string
	Sleep          int `ini:"sleep"`    // url打开失败后等待多久继续 // 秒
	Retry          int `ini:"retry"`    // 失败重试次数
	Maxbody        int `ini:"max_body"` // 允许下载的最大body
	// 表示连接池对所有host的最大链接数量，host也即dest-ip，默认为无穷大（0），
	// 但是通常情况下为了性能考虑都要严格限制该数目（实际使用中通常利用压测二分得到该参数的最佳近似值）。
	// 太大容易导致客户端和服务端的socket数量剧增，导致内存吃满，文件描述符不足等问题；
	// 太小则限制了连接池的socket数量，资源利用率较低。
	MaxIdleConns int `ini:"maxidleconns"`
	// 表示连接池对每个host的最大链接数量，从字面意思也可以看出：
	Maxidleconnsperhost int `ini:"maxidleconnsperhost"`
	// 空闲timeout设置，也即socket在该时间内没有交互则自动关闭连接,该参数通常设置为分钟级别，例如：90秒。
	Idleconntimeout time.Duration `ini:"idleconntimeout"`
	// request header 超时
	ResponseHeaderTimeout time.Duration `ini:"response_header_timeout"`
	ExpectContinueTimeout time.Duration `ini:"expect_continue_timeout"`
	// 确定，这是建立连接的时间
	// 这里主要是控制解析域名时间,TCP 建立连接超时时长设置
	Tcpconnecttimeout time.Duration `ini:"tcpconnecttimeout"`
	// time.Duration不是一个函数，只是将数据显示转为 time.Duration这个类型
	// 这个超时是总的超时时间，如果不清楚Transport里面的设置可以设置这个超时
	// 由这个超时来设置总的超时时长
	Timeouttotal time.Duration `ini:"timeouttotal"`
	Allowredict  bool          `ini:"allowredict"` // 是否允许重定向。默认不允许，使用严格模式来判断
}

// Curlresult
// 请求结果
type Curlresult struct {
	Body       []byte // 相应body
	Header     http.Header
	Url        *url.URL // 请求的URL对象
	HttpStatus int      // http状态
}
