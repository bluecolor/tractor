package tractor

type Output interface {
	PluginDescriber

	Write(wire Wire) error
}
