package server

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type Service interface {
	Serve() error
	CheckHealth()
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

	log.Printf("Creating RoundRobinService. Name: %s, ListenAddr: %s", name, listendAddr)

	return &RoundRobinService{
		Name:            name,
		ListenAddr:      listendAddr,
		RoundRobinCount: 0,
		Endpoints:       endpoints,
	}, nil
}

func (rrs *RoundRobinService) Serve() error {
	log.Printf("Listening on: %s", rrs.ListenAddr)
	return http.ListenAndServe(rrs.ListenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, err := rrs.GetRoute()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		rrs.Endpoints[route].Proxy.ServeHTTP(w, r)
	}))
}

func (rrs *RoundRobinService) GetRoute() (int, error) {
	for i := 0; i < len(rrs.Endpoints); i++ {
		route := rrs.RoundRobinCount
		rrs.RoundRobinCount++
		if rrs.RoundRobinCount >= len(rrs.Endpoints) {
			rrs.RoundRobinCount = 0
		}
		if rrs.Endpoints[route].Healthy {
			return route, nil
		}
	}

	return 0, errors.New("No healthy endpoints")
}

func (rrs *RoundRobinService) CheckHealth() {
	for {
		for _, e := range rrs.Endpoints {
			e.CheckHealth()
		}
		time.Sleep(10 * time.Second)
	}
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

func (lcs *LeastConnectionService) Serve() error {
	log.Print("serving least connection: " + lcs.Name)
	return nil
}

func (lcs *LeastConnectionService) CheckHealth() {
	for {
		for _, e := range lcs.Endpoints {
			e.CheckHealth()
		}
		time.Sleep(10 * time.Second)
	}
}

type IPHashService struct {
	Name       string
	ListenAddr string
	Endpoints  []*Endpoint
}

func NewIPHashService(name string, listendAddr string, endpointStrs []string) (*IPHashService, error) {
	return nil, nil
}

func (ihs *IPHashService) Serve() error {
	log.Print("serving ip hash: " + ihs.Name)
	return nil
}

func (ihs *IPHashService) CheckHealth() {
	for {
		for _, e := range ihs.Endpoints {
			e.CheckHealth()
		}
		time.Sleep(10 * time.Second)
	}
}
