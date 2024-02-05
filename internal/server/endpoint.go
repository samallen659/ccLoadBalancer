package server

import "net/http/httputil"

type Endpoint struct {
	addr  string
	proxy *httputil.ReverseProxy
}
