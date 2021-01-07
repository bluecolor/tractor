package main

import (
	"github.com/bluecolor/tractor/api"
)

type config struct {
	Libdir           string
	Username         string
	Password         string
	ConnectionString string
}

func getConfig(conf api.Config) *config {
	config := config{}
	api.ParseConfig(conf, &config)
	return &config
}

// Run plugin
func Run(conf api.Config) {
	config := getConfig(conf)
	println(config.Username)
}
