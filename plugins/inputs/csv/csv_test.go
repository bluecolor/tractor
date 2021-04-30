package csv

import (
	"testing"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/utils"
	"github.com/stretchr/testify/require"
)

const pluginType string = "csv"

func TestCsv_PluginExists(t *testing.T) {
	registry := inputs.Inputs
	_, ok := registry[pluginType]
	require.True(t, ok)
}

func TestCsv_Create(t *testing.T) {
	creator, _ := inputs.Inputs[pluginType]
	_, err := creator(nil, nil, nil)
	require.NoError(t, err)
}

func TestCsv_Validate(t *testing.T) {
	creator, _ := inputs.Inputs[pluginType]
	params, err := utils.JSONLoadString(`{"path": "/test/path"}`)
	require.NoError(t, err)
	plugin, err := creator(nil, nil, params)
	require.NoError(t, err)
	v, ok := plugin.(tractor.Validator)
	require.True(t, ok)
	err = v.Validate()
	require.NoError(t, err)
}

func TestCsv_ValidateError(t *testing.T) {
	creator, _ := inputs.Inputs[pluginType]
	params, err := utils.JSONLoadString(`{}`)
	require.NoError(t, err)
	plugin, err := creator(nil, nil, params)
	require.NoError(t, err)
	v, ok := plugin.(tractor.Validator)
	require.True(t, ok)
	err = v.Validate()
	require.Error(t, err)
}
