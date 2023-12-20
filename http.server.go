package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HttpServerStart 公功 http server 启动函数
func (h *HttpServer) HttpServerStart() {
	mux := http.NewServeMux()
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
		Handler:           mux,
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
	defer Closehttpserver(server)

	Log("启动Http(s) Server", h.ListenAddr)
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
			caCrt, err := ioutil.ReadFile(Fileabs(h.Ca))
			if err != nil {
				Error("HTTPS Service Load Ca error", err)
				os.Exit(1)
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(caCrt)
			server.TLSConfig.ClientCAs = pool
			server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if err = server.ListenAndServeTLS(Fileabs(h.Crt), Fileabs(h.Key)); err != nil {
			Error("HTTPS Service 服务启动失败", server.Addr, err)
			os.Exit(1)
		}
		return
	}
	go h.hotUpdate(server)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			Error("HTTP Service 启动失败", server.Addr, err)
			os.Exit(1)
		}
	}
	if h.Hotupdate {
		h.HttpServerStart()
	}
}

func (h *HttpServer) middleware(mux *http.ServeMux, u string, f func(w http.ResponseWriter, r *http.Request)) {
	mux.Handle(u, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer CloseReq(r)
		Debug("请求地址", r.URL.String(), "IP", Getip(r))
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
			if Searchslice(addr, h.Allowip) == false {
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
			if Searchslice(addr, h.Denyip) == true {
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
			Error("url 异常", err)
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

// MethodNotAllow 405
func MethodNotAllow(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, `
<h1 style="text-align:center;">400 Error!</h1>
<p style="text-align:center;">`+http.StatusText(http.StatusMethodNotAllowed)+`</p>
`)
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
		Booltostring(this.Ssl),
		this.Ca,
		this.Crt,
		this.Key,
		this.SocketTimeout.String(),
	}, append(this.Allowip, this.Denyip...)...)
	return Md5string(strings.Join(strArr, ""))
}
