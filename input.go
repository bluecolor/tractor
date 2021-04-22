package tractor

type Input interface {
	PluginDescriber

	Read(ch chan<- *Message) error
}

type ServiceInput interface {
	Input

	Start() error

	Stop()
}
