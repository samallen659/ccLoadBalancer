package server

import "log"

type Service interface {
	Serve()
}

type RoundRobinService struct {
	Name            string
	ListenAddr      string
	RoundRobinCount int
	Endpoints       []*Endpoint
}

func NewRoundRobinService(name string, listendAddr string, endpointStrs []string) (*RoundRobinService, error) {
	var endpoints []*Endpoint
	for _, endpointStr := range endpointStrs {
		endpoint, err := NewEndpoint(endpointStr)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	return &RoundRobinService{
		Name:            name,
		ListenAddr:      listendAddr,
		RoundRobinCount: 0,
		Endpoints:       endpoints,
	}, nil
}

func (rrs *RoundRobinService) Serve() {
	log.Print("serving round robin: " + rrs.Name)
}

type LeastConnectionService struct {
	Name            string
	ListenAddr      string
	ConnectionCount map[string]int
	Endpoints       []*Endpoint
}

func NewLeastConnectionService(name string, listendAddr string, endpointStrs []string) (*LeastConnectionService, error) {
	return nil, nil
}

func (lcs *LeastConnectionService) Serve() {
	log.Print("serving least connection: " + lcs.Name)
}

type IPHashService struct {
	Name       string
	ListenAddr string
	Endpoints  []*Endpoint
}

func NewIPHashService(name string, listendAddr string, endpointStrs []string) (*IPHashService, error) {
	return nil, nil
}

func (ihs *IPHashService) Serve() {
	log.Print("serving ip hash: " + ihs.Name)
}
