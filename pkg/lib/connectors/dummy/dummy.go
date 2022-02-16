package dummy

import (
	"fmt"
	"reflect"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"gorm.io/gorm/utils"
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

func (c *DummyConnector) Validate(config connectors.ConnectorConfig) error {
	fields := reflect.VisibleFields(reflect.TypeOf(c.config))
	tags := make([]string, len(fields))
	for i, field := range fields {
		tags[i] = field.Tag.Get("json")
	}
	for key, _ := range config {
		if !utils.Contains(tags, key) {
			return fmt.Errorf("invalid config key %s", key)
		}
	}
	return nil
}

func init() {
	connectors.Add("dummy", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config)
	})
}
