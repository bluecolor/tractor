package util

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/bluecolor/tractor/api"
	"gopkg.in/yaml.v2"
)

// Mapping self contained source and target conf
type Mapping struct {
	Input  api.Config
	Output api.Config
}

// GetMappings get mappings in config file
func GetMappings(configFile string) []map[string]Mapping {
	var mappings []map[string]Mapping
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &mappings)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return mappings
}

// GetMapping Get mapping with given name
func GetMapping(configFile string, name string) (*Mapping, error) {
	mappings := GetMappings(configFile)
	for _, mapping := range mappings {
		if m, ok := mapping[name]; ok {
			return &m, nil
		}
	}
	return nil, errors.New("Mapping with name: " + name + " not found")
}
