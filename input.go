package tractor

type Input interface {
	PluginDescriber

	Read(wire Wire) error
}

type ServiceInput interface {
	Input

	Start() error

	Stop()
}
