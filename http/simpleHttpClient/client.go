package simpleHttpClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

var simpleClient *http.Client

// InitHttpClient 初始化http client
func InitHttpClient(timeout time.Duration, args ...string) (*http.Client, error) {
	if simpleClient != nil {
		return simpleClient, nil
	}

	trans := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // client 不对https 证书进行校验
		},
		DisableKeepAlives:   false,           // 是否禁用 keep-alives
		TLSHandshakeTimeout: 5 * time.Second, // tls 握手超时
		DisableCompression:  false,           // 是否禁用压缩，默认情况下应保持 false 以允许接收压缩内容，除非有特定原因需要禁用。
		// 连接池相关配置
		MaxIdleConns:        1000, // 控制整个客户端的最大空闲连接数
		MaxIdleConnsPerHost: 5,    // 限制每个主机的最大空闲连接数
		MaxConnsPerHost:     0,    // 限制每个主机的最大连接数（包括活跃和空闲）
		IdleConnTimeout:     0,    // 设置空闲连接在被关闭前等待新请求的时间
	}
	
	if len(args) >= 2 {
		proxyAddr := args[1]
		switch args[0] {
		case "socks5":
			dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
				return simpleClient, fmt.Errorf("newHttpClient socks5 proxy error: %v", err)
			}
			trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		case "http":
			u, err := url.Parse(proxyAddr)
			if err != nil {
				return simpleClient, fmt.Errorf("newHttpClient parse proxy url error: %v", err)
			}
			trans.Proxy = http.ProxyURL(u)
		}
	}
	simpleClient = &http.Client{
		Transport: trans,
		Timeout:   timeout,
	}
	return simpleClient, nil
}
