package providers

import "encoding/json"

type ProviderConfig []byte

func (c *ProviderConfig) LoadConfig(config interface{}) error {
	return json.Unmarshal([]byte(*c), config)
}
