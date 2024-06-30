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
var impatientHTTPClient *http.Client

func init() {
	transport := &http.Transport{
		Proxy:             ProxyFunc,
		DialContext:       DialContextFunc,
		DisableKeepAlives: true,
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

	impatientHTTPClient = &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}

func GetHttpClient() *http.Client {
	return httpClient
}

func GetImpatientHttpClient() *http.Client {
	return impatientHTTPClient
}

func ProxyFunc(req *http.Request) (*url.URL, error) {
	proxy, ok := (req.Context().Value("proxy")).(string)
	if !ok || proxy == "" || (!strings.HasPrefix(proxy, "http://") && !strings.HasPrefix(proxy, "https://")) {
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

func DialContextFunc(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	proxy, ok := ctx.Value("proxy").(string)
	if !ok || proxy == "" || !strings.HasPrefix(proxy, "socks5://") {
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
