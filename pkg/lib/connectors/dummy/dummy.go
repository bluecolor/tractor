package dummy

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors"
)

type DummyConfig struct {
	GenerateFakeData   bool `json:"generateFakeData"`   // self generate fake data
	FakeRecordCount    int  `json:"fakeRecordCount"`    // total fake record count to generate
	FakeRecordInterval int  `json:"fakeRecordInterval"` // sleep per send in milliseconds
}

type DummyConnector struct {
	connectors.Base
	config DummyConfig
}

const (
	InputChannelKey         = "channel"
	OutputChannelKey        = "channel"
	InputMessageChannelKey  = "message_channel"
	OutputMessageChannelKey = "message_channel"
)

func New(config connectors.ConnectorConfig) (*DummyConnector, error) {
	mysqlConfig := DummyConfig{}
	if err := config.LoadConfig(&mysqlConfig); err != nil {
		return nil, err
	}
	return &DummyConnector{
		config: mysqlConfig,
	}, nil
}

func init() {
	connectors.Add("dummy", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config)
	})
}
