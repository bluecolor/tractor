package session

import (
	"errors"

	"github.com/bluecolor/tractor/lib/plugins"
)

func (s *Session) ValidateInputPlugin() error {
	if validator, ok := s.InputPlugin.(plugins.Validator); ok {
		if err := validator.Validate(); err != nil {
			return errors.New("❌  Failed to validate plugin config")
		} else {
			return nil
		}
	}
	return nil
}

func (s *Session) ValidateOutputPlugin() error {
	if validator, ok := s.OutputPlugin.(plugins.Validator); ok {
		if err := validator.Validate(); err != nil {
			return errors.New("❌  Failed to validate plugin config")
		} else {
			return nil
		}
	}
	return nil
}
