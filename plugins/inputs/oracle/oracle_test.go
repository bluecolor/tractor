package oracle

import (
	"testing"

	"github.com/stretchr/stew/slice"
	"github.com/stretchr/testify/require"

	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
)

func getPlugin(path string) (interface{}, error) {
	registry := inputs.Inputs
	c := config.NewConfig()
	c.LoadConfig(path)
	return registry["oracle"](c.Mappings[0].Input.Config, nil, nil)
}

func TestOracle_CheckRegistry(t *testing.T) {
	registry := inputs.Inputs
	keys := []string{}

	for key := range registry {
		keys = append(keys, key)
	}
	require.Equal(t, true, slice.Contains(keys, "oracle"))
}

func TestOracle_LoadConfig(t *testing.T) {
	plugin, err := getPlugin("oracle_test.yml")
	oracle := plugin.(*Oracle)
	require.NoError(t, err)
	require.Equal(t, "tractor", oracle.Username)
	require.Equal(t, "tractor", oracle.Password)
	require.Equal(t, 1521, oracle.Port)
	require.Equal(t, 1, oracle.Parallel)
	require.Equal(t, 1000, oracle.FetchSize)
}

func TestOracle_Validate(t *testing.T) {
	plugin, err := getPlugin("oracle_test.yml")
	require.NoError(t, err)
	oracle, ok := plugin.(*Oracle)
	require.True(t, ok)
	require.NoError(t, oracle.Validate())
}

func TestOracle_Initializer(t *testing.T) {
	// requires connection
	// plugin, err := getPlugin("oracle_test.yml")
	// require.NoError(t, err)
	// if i, ok := plugin.(tractor.Initializer); ok {
	// 	require.NoError(t, i.Init())
	// }
}
