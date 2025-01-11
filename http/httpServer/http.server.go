package httpServer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gorilla/handlers"
	"github.com/helays/utils/close/httpClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/crypto/md5"
	"github.com/helays/utils/logger/ulogs"
	"github.com/helays/utils/logger/zaploger"
	"github.com/helays/utils/tools"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HttpServerStart 公功 http server 启动函数
func (h *HttpServer) HttpServerStart() {
	h.serverNameMap = make(map[string]byte)
	for _, dom := range h.ServerName {
		h.serverNameMap[strings.ToLower(dom)] = 0
	}
	mux := http.NewServeMux()
	h.Route["/debug/switch-debug"] = SwitchDebug
	if config.Dbg {
		h.Route["/debug/pprof/"] = pprof.Index
		h.Route["/debug/pprof/cmdline"] = pprof.Cmdline
		h.Route["/debug/pprof/profile"] = pprof.Profile
		h.Route["/debug/pprof/symbol"] = pprof.Symbol
		h.Route["/debug/pprof/trace"] = pprof.Trace
	}
	if h.Logger.LogLevelConfigs != nil {
		h.logger, err = zaploger.New(&h.Logger)
		ulogs.DieCheckerr(err, "http server 日志模块初始化失败")
	}

	if h.Route != nil {
		for u, funcName := range h.Route {
			h.middleware(mux, u, funcName)
		}
	}
	if h.RouteSocket != nil {
		for u, funcName := range h.RouteSocket {
			mux.Handle(u, websocket.Handler(funcName))
		}
	}

	server := &http.Server{
		Addr:              h.ListenAddr,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		// BaseContext:       nil,
		// ConnContext:       nil,
	}
	if h.EnableGzip {
		server.Handler = handlers.CompressHandler(mux)
	} else {
		server.Handler = mux
	}
	defer Closehttpserver(server)

	ulogs.Log("启动Http(s) Server", h.ListenAddr)
	if h.Ssl {
		server.TLSConfig = &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		}
		// 如果包含ca证书，就需要做强制双向https 验证
		if h.Ca != "" {
			caCrt, err := os.ReadFile(tools.Fileabs(h.Ca))
			if err != nil {
				ulogs.Error("HTTPS Service Load Ca error", err)
				os.Exit(1)
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(caCrt)
			server.TLSConfig.ClientCAs = pool
			server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if err = server.ListenAndServeTLS(tools.Fileabs(h.Crt), tools.Fileabs(h.Key)); err != nil {
			ulogs.Error("HTTPS Service 服务启动失败", server.Addr, err)
			os.Exit(1)
		}
		return
	}
	go h.hotUpdate(server)
	var isQuit bool
	go h.stopServer(server, &isQuit)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			ulogs.Error("HTTP Service 启动失败", server.Addr, err)
			os.Exit(1)
		}
	}
	if h.Hotupdate && !isQuit {
		h.HttpServerStart()
	}
}

func (h *HttpServer) middleware(mux *http.ServeMux, u string, f func(w http.ResponseWriter, r *http.Request)) {
	mux.Handle(u, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer httpClose.CloseReq(r)
		if len(h.serverNameMap) > 0 {
			// 需要验证请求的host是否绑定了
			_host := strings.Split(strings.ToLower(r.Host), ":")
			if len(_host) > 0 {
				if _, ok := h.serverNameMap[_host[0]]; !ok {
					w.WriteHeader(http.StatusBadGateway)
					return
				}
			}
		}
		if h.logger != nil {
			// 这里输出info 级别的请求日志
			h.logger.Info(context.Background(),
				Getip(r),
				zap.String(r.Method, r.URL.String()),
				zap.String("http_user_agent", r.Header.Get("User-Agent")),
			)
		} else {
			ulogs.Debug("请求地址", r.URL.String(), "IP", Getip(r))
		}

		// add header
		w.Header().Set("server", "vs/1.0")
		w.Header().Set("connection", "keep-alive")
		// 白名单验证
		if len(h.Allowip) > 0 {
			// 存在白名单，只允许白
			//名单中存在的访问
			addr := r.RemoteAddr
			al := strings.Index(addr, ":")
			if al < 0 {
				Forbidden(w, "")
				return
			}
			addr = addr[0:al]
			if tools.Searchslice(addr, h.Allowip) == false {
				Forbidden(w, "")
				return
			}
		} else if len(h.Denyip) > 0 {
			// 存在黑名单的话，再黑名单中的 IP禁止访问
			addr := r.RemoteAddr
			al := strings.Index(addr, ":")
			if al < 0 {
				Forbidden(w, "")
				return
			}
			addr = addr[0:al]
			if tools.Searchslice(addr, h.Denyip) == true {
				Forbidden(w, "")
				return
			}
		}

		if h.CommonCallback != nil && !h.CommonCallback(w, r) {
			return
		}
		http.HandlerFunc(f).ServeHTTP(w, r)
	}))
}

// SetRequestDefaultPage 设置 打开的默认页面
// defaultPage string 默认打开页面
// root 网站更目录
// path string
func SetRequestDefaultPage(defaultPage, path string) ([]*os.File, []string, bool) {
	sarr := strings.Split(path, "??")
	if len(sarr) == 1 {
		swapUrl, err := url.Parse(path)
		if err != nil {
			ulogs.Error("url 异常", err)
			return nil, nil, false
		}
		path = swapUrl.Path
		if filepath.Base(path) == "lib.js" {

		}
		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return []*os.File{f}, []string{path}, false
		}
		fInfo, _ := f.Stat()
		if !fInfo.IsDir() {
			return []*os.File{f}, []string{path}, true
		}
		defaultPage = strings.TrimSpace(defaultPage)
		if defaultPage == "" {
			defaultPage = "index.html"
		}
		fp := path + "/" + defaultPage

		if strings.HasSuffix(path, "/") {
			fp = path + defaultPage
		}
		f, err = os.OpenFile(fp, os.O_RDONLY, 0644)
		if err != nil {
			return []*os.File{f}, []string{fp}, false
		}
		return []*os.File{f}, []string{fp}, true
	}

	var (
		swapList  []*os.File
		swapPaths []string
		status    bool
	)
	for _, item := range strings.Split(sarr[1], ",") {
		swapFile, swapPath, swapStatus := SetRequestDefaultPage(defaultPage, sarr[0]+item)
		if !swapStatus {
			continue
		}
		swapList = append(swapList, swapFile...)
		swapPaths = append(swapPaths, swapPath...)
	}
	if len(swapList) > 0 {
		status = true
	}
	return swapList, swapPaths, status
}

// Closehttpserver 关闭http server
func Closehttpserver(s *http.Server) {
	if s != nil {
		_ = s.Close()
	}
}

// 用于检测参数变更，然后热更新。
func (this *HttpServer) hotUpdate(server *http.Server) {
	if this.Hotupdate {
		go func() {
			hash := this.hash()
			tck := time.NewTicker(1 * time.Second)
			defer tck.Stop()
			for range tck.C {
				if hash == this.hash() {
					continue
				}
				tck.Stop()
				break
			}
			Closehttpserver(server)
		}()
	}
}

// 计算 httpserver 模块摘要
func (this HttpServer) hash() string {
	strArr := append([]string{
		this.ListenAddr,
		this.Auth,
		tools.Booltostring(this.Ssl),
		this.Ca,
		this.Crt,
		this.Key,
		this.SocketTimeout.String(),
	}, append(this.Allowip, this.Denyip...)...)
	return md5.Md5string(strings.Join(strArr, ""))
}

func (this HttpServer) stopServer(server *http.Server, isQuit *bool) {
	config.EnableHttpserver = true
	_ = <-config.CloseHttpserverSig
	*isQuit = true
	ulogs.Log("http server已关闭")
	Closehttpserver(server)
	config.CloseHttpserverSig <- 1
}
