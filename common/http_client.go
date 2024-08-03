package common

import (
	"context"
	"errors"
	p "golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetImageHttpClient() (*http.Client, error) {
	return GetProxiedHttpClient(OutProxyUrl)
}

func GetProxiedHttpClient(proxyUrl string) (*http.Client, error) {
	if proxyUrl == "" {
		return &http.Client{}, nil
	}
	u, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(proxyUrl, "http://") || strings.HasPrefix(proxyUrl, "https://"):
		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}, nil
	case strings.HasPrefix(proxyUrl, "socks5://"):
		dialer, err := p.FromURL(u, p.Direct)
		if err != nil {
			return nil, err
		}
		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.(p.ContextDialer).DialContext(ctx, network, addr)
				},
			},
		}, nil
	default:
		return nil, errors.New("unsupported proxy type: " + proxyUrl)
	}
}
