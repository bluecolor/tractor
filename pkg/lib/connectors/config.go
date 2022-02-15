package connectors

import (
	"github.com/bluecolor/tractor/pkg/utils"
)

type ConfigValidator interface {
	Validate(config ConnectorConfig) error
}

type ConnectorConfig map[string]interface{}

func (c ConnectorConfig) LoadConfig(config interface{}) error {
	if err := utils.MapToStruct(c, config); err != nil {
		return err
	}
	return nil
}
