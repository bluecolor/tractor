package util

import (
	"errors"
	"io/ioutil"

	"github.com/bluecolor/tractor/api"
	"gopkg.in/yaml.v2"
)

// Mapping self contained source and target conf
type Mapping struct {
	Input  api.Config
	Output api.Config
}

// GetMappings get mappings in config file
func GetMappings(mappingsFile string) ([]map[string]Mapping, error) {
	var mappings []map[string]Mapping
	yamlFile, err := ioutil.ReadFile(mappingsFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

// ValidateMapping ...
func ValidateMapping(mapping *Mapping) error {
	if _, ok := mapping.Input["plugin"]; !ok {
		return errors.New("Input plugin type is not given")
	}

	if _, ok := mapping.Output["plugin"]; !ok {
		return errors.New("Input plugin type is not given")
	}

	return nil
}

// GetMapping Get mapping with given name
func GetMapping(mappingsFile string, name string, args ...interface{}) (*Mapping, error) {
	mappings, err := GetMappings(mappingsFile)
	if err != nil {
		return nil, err
	}
	for _, mapping := range mappings {
		if m, ok := mapping[name]; ok {
			if len(args) > 0 && args[0].(bool) {
				if err := ValidateMapping(&m); err != nil {
					return nil, err
				}
			}
			return &m, nil
		}
	}
	return nil, errors.New("Mapping with name: " + name + " not found")
}

// GetConfigs ...
func GetConfigs(m *Mapping) ([]byte, []byte, error) {
	iconf, err := yaml.Marshal(m.Input)
	if err != nil {
		return nil, nil, err
	}
	oconf, err := yaml.Marshal(m.Output)
	if err != nil {
		return nil, nil, err
	}
	return iconf, oconf, nil
}
