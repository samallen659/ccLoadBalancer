package server

import (
	"errors"
	"hash/fnv"
	"log"
	"math"
	"net/http"
	"time"
)

var noHealthyEndpointsError = errors.New("No healthy endpoints")

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
		rrs.Endpoints[route].IncrementConnection()
		rrs.Endpoints[route].Proxy.ServeHTTP(w, r)
		rrs.Endpoints[route].DecrementConnection()
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

	return 0, noHealthyEndpointsError
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
	Name       string
	ListenAddr string
	Endpoints  []*Endpoint
}

func NewLeastConnectionService(name string, listendAddr string, endpointStrs []string) (*LeastConnectionService, error) {
	var endpoints []*Endpoint
	for _, endpointStr := range endpointStrs {
		endpoint, err := NewEndpoint(endpointStr)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	log.Printf("Creating LeastConnectionService. Name: %s, ListenAddr: %s", name, listendAddr)

	return &LeastConnectionService{
		Name:       name,
		ListenAddr: listendAddr,
		Endpoints:  endpoints,
	}, nil
}

func (lcs *LeastConnectionService) Serve() error {
	log.Printf("Listening on: %s", lcs.ListenAddr)
	return http.ListenAndServe(lcs.ListenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, err := lcs.GetRoute()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		lcs.Endpoints[route].IncrementConnection()
		lcs.Endpoints[route].Proxy.ServeHTTP(w, r)
		lcs.Endpoints[route].DecrementConnection()
	}))
}

func (lcs *LeastConnectionService) GetRoute() (int, error) {
	var route int
	connectionCount := math.MaxInt
	for i, e := range lcs.Endpoints {
		if e.ConnectionCount < connectionCount {
			route = i
		}
	}

	if lcs.Endpoints[route].Healthy {
		return route, nil
	}

	for i := 0; i < len(lcs.Endpoints)-1; i++ {
		route++
		if route >= len(lcs.Endpoints) {
			route = 0
		}
		if lcs.Endpoints[route].Healthy {
			return route, nil
		}
	}

	return 0, noHealthyEndpointsError
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
	var endpoints []*Endpoint
	connectionCount := make(map[string]int)
	for _, endpointStr := range endpointStrs {
		connectionCount[endpointStr] = 0
		endpoint, err := NewEndpoint(endpointStr)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	log.Printf("Creating IPHashService. Name: %s, ListenAddr: %s", name, listendAddr)

	return &IPHashService{
		Name:       name,
		ListenAddr: listendAddr,
		Endpoints:  endpoints,
	}, nil
}

func (ihs *IPHashService) Serve() error {
	log.Print("Listening on: %s" + ihs.Name)
	return http.ListenAndServe(ihs.ListenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, err := ihs.GetRoute(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		ihs.Endpoints[route].IncrementConnection()
		ihs.Endpoints[route].Proxy.ServeHTTP(w, r)
		ihs.Endpoints[route].DecrementConnection()
	}))
}

func (ihs *IPHashService) CheckHealth() {
	for {
		for _, e := range ihs.Endpoints {
			e.CheckHealth()
		}
		time.Sleep(10 * time.Second)
	}
}

func (ihs *IPHashService) GetRoute(addr string) (int, error) {
	h := fnv.New32a()
	h.Write([]byte(addr))
	hash := h.Sum32()
	route := int(hash) % len(ihs.Endpoints)

	if ihs.Endpoints[route].Healthy {
		return route, nil
	}

	for i := 0; i < len(ihs.Endpoints)-1; i++ {
		route++
		if route >= len(ihs.Endpoints) {
			route = 0
		}
		if ihs.Endpoints[route].Healthy {
			return route, nil
		}
	}

	return 0, noHealthyEndpointsError
}
