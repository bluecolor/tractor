package api

import (
	"gopkg.in/yaml.v2"
)

// ParseConfig parse and load given conf to config struct
func ParseConfig(conf Config, config interface{}) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, config)
	return nil
}
