package server

import (
	"net/http/httputil"
	"net/url"
)

type Endpoint struct {
	Addr    *url.URL
	Proxy   *httputil.ReverseProxy
	Healthy bool
}

func NewEndpoint(addr string) (*Endpoint, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return &Endpoint{
		Addr:    url,
		Proxy:   proxy,
		Healthy: false,
	}, nil
}
