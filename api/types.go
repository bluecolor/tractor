package api

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
