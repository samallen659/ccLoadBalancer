package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Endpoint struct {
	Addr            *url.URL
	Proxy           *httputil.ReverseProxy
	Healthy         bool
	ConnectionCount int
	mu              sync.Mutex
}

func NewEndpoint(addr string) (*Endpoint, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	log.Printf("Creating Endpoint with url: %s", url.String())

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = url.String()
			req.Host = url.String()
		},
	}

	return &Endpoint{
		Addr:            url,
		Proxy:           proxy,
		Healthy:         false,
		ConnectionCount: 0,
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

func (e *Endpoint) IncrementConnection() {
	e.mu.Lock()
	e.ConnectionCount++
	e.mu.Unlock()
}

func (e *Endpoint) DecrementConnection() {
	e.mu.Lock()
	e.ConnectionCount--
	e.mu.Unlock()
}
