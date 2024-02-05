package server

import "log"

type Service interface {
	Serve()
}

type RoundRobinService struct {
	Name            string
	ListenAddr      string
	RoundRobinCount int
	Endpoints       []Endpoint
}

func (rrs *RoundRobinService) Serve() {
	log.Print("serving round robin: " + rrs.Name)
}

type LeastConnectionService struct {
	Name            string
	ListenAddr      string
	ConnectionCount map[string]int
	Endpoints       []Endpoint
}

func (lcs *LeastConnectionService) Serve() {
	log.Print("serving least connection: " + lcs.Name)
}

type IPHashService struct {
	Name       string
	ListenAddr string
	Endpoints  []Endpoint
}

func (ihs *IPHashService) Serve() {
	log.Print("serving ip hash: " + ihs.Name)
}
