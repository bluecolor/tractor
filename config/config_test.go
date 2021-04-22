package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmd_LoadConfig(t *testing.T) {
	c := NewConfig()
	require.NoError(t, c.LoadConfig("../test/tractor_test.yml"))
	require.Equal(t, "Hello", c.Mappings[0].Name)
}
