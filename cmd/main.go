package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Services []Service `yaml:"config"`
}

type Service struct {
	Name       string   `yaml:"name"`
	ListenAddr string   `yaml:"listenAddr"`
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

	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Panic(err)
	}

	log.Print(c.Services[0])
	log.Print(c.Services[1])
}
