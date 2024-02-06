package main

import (
	"fmt"
	"github.com/samallen659/ccLoadBalancer/internal/server"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/config_example.yaml", dir))
	if err != nil {
		log.Panic(err)
	}

	var c server.Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(c.Services)

	s, err := server.NewServer(c)
	if err != nil {
		log.Panic(err)
	}

	if err = s.Serve(); err != nil {
		log.Panic(err)
	}
}
