package server

import "net/http/httputil"

type Endpoint struct {
	Addr    string
	Proxy   *httputil.ReverseProxy
	Healthy bool
}
