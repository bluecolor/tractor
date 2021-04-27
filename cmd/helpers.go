package cmd

import (
	"errors"
	"fmt"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/plugins/outputs"
)

func validateAndGetInputPlugin(name string, config map[string]interface{}) (tractor.Input, error) {
	if creator, ok := inputs.Inputs[name]; ok {
		inputPlugin := creator(config)
		if validator, ok := inputPlugin.(tractor.Validator); ok {
			if err := validator.ValidateConfig(); err != nil {
				return nil, errors.New("❌  Failed to validate input config")
			} else {
				println("☑️  Input config validated")
				return inputPlugin, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("❌  No matching input plugin found for %s", name))
}

func validateAndGetOutputPlugin(name string, config map[string]interface{}) (tractor.Output, error) {
	if creator, ok := outputs.Ouputs[name]; ok {
		outputPlugin := creator(config)
		if validator, ok := outputPlugin.(tractor.Validator); ok {
			if err := validator.ValidateConfig(); err != nil {
				return nil, errors.New("❌  Failed to validate output config")
			} else {
				println("☑️  Output config validated")
				return outputPlugin, nil
			}
		} else {
			println("Validator not implemented for output plugin")
			return outputPlugin, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("❌  No matching output plugin found for %s", name))
}
