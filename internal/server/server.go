package server

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Config struct {
	Services []struct {
		Name       string   `yaml:"name"`
		ListenAddr string   `yaml:"listenAddr"`
		Algorithm  string   `yaml:"algorithm"`
		Endpoints  []string `yaml:"endpoints"`
	} `yaml:"config"`
}

type Server struct {
	Services []Service
}

func NewServer(config Config) (*Server, error) {
	var services []Service
	for _, s := range config.Services {
		switch s.Algorithm {
		case "roundrobin":
			sv, err := NewRoundRobinService(s.Name, s.ListenAddr, s.Endpoints)
			if err != nil {
				return nil, err
			}
			services = append(services, sv)
		case "leastconnection":
			sv, err := NewLeastConnectionService(s.Name, s.ListenAddr, s.Endpoints)
			if err != nil {
				return nil, err
			}
			services = append(services, sv)
		case "iphash":
			sv, err := NewIPHashService(s.Name, s.ListenAddr, s.Endpoints)
			if err != nil {
				return nil, err
			}
			services = append(services, sv)
		default:
			return nil, fmt.Errorf("Unknown routing algorithm: %s", s.Algorithm)
		}
	}

	return &Server{services}, nil
}

func (s *Server) Serve() error {
	eg := errgroup.Group{}
	for _, sv := range s.Services {
		go sv.CheckHealth()
		eg.Go(sv.Serve)
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
