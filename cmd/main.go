package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Name       string   `yaml:"name"`
	ListenAddr net.Addr `yaml:"listenAddr"`
	Algorithm  string   `yaml:"algorithm"`
	Endpoints  []string `yaml:"endpoints"`
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/config_example.yaml", dir))
	if err != nil {
		log.Panic(err)
	}

	var s Service
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Panic(err)
	}

	log.Print(s)
}
