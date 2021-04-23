package tractor

type Output interface {
	PluginDescriber

	Write(ch <-chan *Message) error
}
