package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c := NewConfig()
	err := c.LoadConfig("../../config/examples/test.yml")
	if err != nil {
		t.Errorf("Error loading config: %s", err)
	}
}
