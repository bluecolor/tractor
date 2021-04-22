package tractor

type Output interface {
	PluginDescriber

	Write(data Data) error
}
