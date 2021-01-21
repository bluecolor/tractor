package api

import (
	"plugin"
	"sync"

	"github.com/bluecolor/tractor/api/message"
)

//Config either input our ooutput configuration given by the user
//in mappings.yml file
type Config map[interface{}]interface{}

// PluginType ...
type PluginType int

const (
	// InputPlugin ...
	InputPlugin PluginType = iota
	// OutputPlugin ...
	OutputPlugin
)

// TractorPlugin ...
type TractorPlugin struct {
	Plugin *plugin.Plugin
	Run    func(*sync.WaitGroup, []byte, chan *message.Message)
}
