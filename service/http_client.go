package service

import (
	"context"
	"fmt"
	p "golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"one-api/common"
	"strings"
	"time"
)

var httpClient *http.Client
var impatientHttpClient *http.Client
var imageHttpClient *http.Client

func init() {
	transport := &http.Transport{
		DisableKeepAlives: true,
		Proxy: ProxyFunc(func(req *http.Request) string {
			proxy, ok := (req.Context().Value("proxy")).(string)
			if !ok {
				return ""
			}
			return proxy
		}),
		DialContext: DialContextFunc(func(ctx context.Context) string {
			proxy, ok := ctx.Value("proxy").(string)
			if !ok {
				return ""
			}
			return proxy
		}),
	}
	if common.RelayTimeout == 0 {
		httpClient = &http.Client{
			Transport: transport,
		}
	} else {
		httpClient = &http.Client{
			Timeout:   time.Duration(common.RelayTimeout) * time.Second,
			Transport: transport,
		}
	}

	impatientHttpClient = &http.Client{
		Timeout:   5 * time.Second, /**/
		Transport: transport,
	}

	// 图片出口代理
	imageHttpClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy: ProxyFunc(func(req *http.Request) string {
				return common.OutProxyUrl
			}),
			DialContext: DialContextFunc(func(ctx context.Context) string {
				return common.OutProxyUrl
			}),
		},
	}
}

func GetHttpClient() *http.Client {
	return httpClient
}

func GetImpatientHttpClient() *http.Client {
	return impatientHttpClient
}

func GetImageHttpClient() *http.Client {
	return imageHttpClient
}

func ProxyFunc(fn func(req *http.Request) string) func(r *http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		proxy := fn(req)
		if proxy == "" || (!strings.HasPrefix(proxy, "http://") && !strings.HasPrefix(proxy, "https://")) {
			return nil, nil
		}
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy address: %w", err)
		}
		switch proxyURL.Scheme {
		case "http", "https":
			return proxyURL, nil
		}
		return nil, fmt.Errorf("unsupported proxy scheme: %s", proxyURL.Scheme)
	}
}

func DialContextFunc(fn func(ctx context.Context) string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		proxy := fn(ctx)
		if proxy == "" || !strings.HasPrefix(proxy, "socks5://") {
			return dialer.DialContext(ctx, network, addr)
		}
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy address: %w", err)
		}
		proxyDialer, err := p.FromURL(proxyURL, dialer)
		if err != nil {
			return nil, fmt.Errorf("error creating proxy dialer: %w", err)
		}
		return proxyDialer.Dial(network, addr)
	}
}
