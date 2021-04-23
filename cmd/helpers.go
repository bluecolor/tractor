package cmd

import (
	"errors"
	"fmt"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/plugins/inputs"
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
	return nil, errors.New(fmt.Sprintf("❌  No matching plugin found for %s", name))
}
