package server

import (
	"net/http"
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

func (e *Endpoint) CheckHealth() {
	resp, err := http.Get(e.Addr.Host)
	if err != nil {
		e.Healthy = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		e.Healthy = false
		return
	}

	e.Healthy = true
}
