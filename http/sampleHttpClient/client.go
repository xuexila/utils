package sampleHttpClient

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

// InitHttpClient 初始化http client
func InitHttpClient(timeout time.Duration, args ...string) (http.Client, error) {
	trans := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // client 不对https 证书进行校验
		},
	}
	if len(args) >= 2 {
		proxyAddr := args[1]
		switch args[0] {
		case "socks5":
			dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
				return http.Client{}, fmt.Errorf("newHttpClient socks5 proxy error: %v", err)
			}
			trans.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		case "http":
			u, err := url.Parse(proxyAddr)
			if err != nil {
				return http.Client{}, fmt.Errorf("newHttpClient parse proxy url error: %v", err)
			}
			trans.Proxy = http.ProxyURL(u)
		}
	}
	return http.Client{
		Transport: trans,
		Timeout:   timeout,
	}, nil
}
