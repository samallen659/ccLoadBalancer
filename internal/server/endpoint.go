package server

import (
	"log"
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
	log.Printf("Creating Endpoint with url: %s", url.String())

	proxy := httputil.NewSingleHostReverseProxy(url)

	return &Endpoint{
		Addr:    url,
		Proxy:   proxy,
		Healthy: false,
	}, nil
}

func (e *Endpoint) CheckHealth() {
	log.Printf("Checking health for: %s", e.Addr.String())
	resp, err := http.Get("http://" + e.Addr.String())
	if err != nil {
		log.Printf("Health check failed for %s: %v", e.Addr.String(), err)
		e.Healthy = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed for %s: %v", e.Addr.String(), err)
		e.Healthy = false
		return
	}

	log.Printf("Health check passed for: %s", e.Addr.String())
	e.Healthy = true
}
