package api

import (
	"gopkg.in/yaml.v2"
)

// ParseConfig parse and load given conf to config struct
func ParseConfig(data []byte, config interface{}) error {
	return yaml.Unmarshal(data, config)
}
