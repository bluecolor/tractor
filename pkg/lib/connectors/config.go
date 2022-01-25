package connectors

import "encoding/json"

type ConnectorConfig []byte

func (c *ConnectorConfig) LoadConfig(config interface{}) error {
	return json.Unmarshal([]byte(*c), config)
}
