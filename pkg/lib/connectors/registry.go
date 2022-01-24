package connectors

import "fmt"

type Creator func(config ConnectorConfig) (Connector, error)

var Connectors = map[string]Creator{}

func Add(name string, creator Creator) {
	Connectors[name] = creator
}

func GetConnector(name string, config ConnectorConfig) (Connector, error) {
	if creator, ok := Connectors[name]; ok {
		return creator(config)
	}
	return nil, fmt.Errorf("Connector %s not found", name)
}
