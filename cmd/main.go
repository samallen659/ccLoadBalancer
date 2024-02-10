package main

import (
	"fmt"
	"github.com/samallen659/ccLoadBalancer/internal/server"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var configFileFlag string

func main() {
	configFliePath := os.Args[1]
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, configFliePath))
	if err != nil {
		log.Panic(err)
	}

	var c server.Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Panic(err)
	}

	s, err := server.NewServer(c)
	if err != nil {
		log.Panic(err)
	}

	if err = s.Serve(); err != nil {
		log.Panic(err)
	}
}
