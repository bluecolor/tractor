package main

import (
	"sync"

	"github.com/bluecolor/tractor/api"
)

type config struct {
	Libdir           string `yaml:"libdir"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	ConnectionString string `yaml:"connection_string"`
}

func getConfig(conf []byte) (*config, error) {
	config := config{}
	if err := api.ParseConfig(conf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Run plugin
func Run(wg *sync.WaitGroup, conf []byte, channel chan []byte) {
	config, err := getConfig(conf)
	if err != nil {
		panic(err)
	}
	println(config.Username)

	var recievedCount int = 0
	for message := range channel {
		recievedCount = recievedCount + 1

		println(string(message))
	}

	wg.Done()
}
