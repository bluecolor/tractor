package file

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/providers"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/file/providers/all"
)

type FileConfig struct {
	ProviderType string                 `json:"providerType"`
	Provider     map[string]interface{} `json:"provider"`
}

type FileConnector struct {
	config   FileConfig
	provider providers.Provider
}

func NewFileConnector(config connectors.ConnectorConfig) (*FileConnector, error) {
	fileConfig := FileConfig{}
	if err := config.LoadConfig(&fileConfig); err != nil {
		return nil, err
	}
	provider, err := providers.GetProvider(fileConfig.ProviderType, fileConfig.Provider)
	if err != nil {
		return nil, err
	}
	return &FileConnector{
		config:   fileConfig,
		provider: provider,
	}, nil
}
