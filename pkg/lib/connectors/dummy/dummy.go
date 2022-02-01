package dummy

import "github.com/bluecolor/tractor/pkg/lib/connectors"

type DummyConnector struct {
	connectors.BaseConnector
}

const (
	InputChannelKey  = "input_channel"
	OutputChannelKey = "output_channel"
)

func New(config connectors.ConnectorConfig) *DummyConnector {
	return &DummyConnector{}
}

func init() {
	connectors.Add("dummy", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config), nil
	})
}
