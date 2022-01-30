package file

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/formats"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/file/formats/all"
	"go.beyondstorage.io/v5/types"
)

type StorageConfig map[string]interface{}

func (c StorageConfig) WithURL(url string) StorageConfig {
	c["url"] = url
	return c
}
func (c StorageConfig) GetURL() string {
	return c["url"].(string)
}

type FileConfig struct {
	StorageType   string        `json:"storageType"`
	Format        string        `json:"format"`
	StorageConfig StorageConfig `json:"storageConfig"`
}

type FileConnector struct {
	Config     FileConfig
	FileFormat *formats.FileFormat
	Storage    types.Storager
}

func New(config connectors.ConnectorConfig) (*FileConnector, error) {
	fc := FileConfig{}
	if err := config.LoadConfig(&fc); err != nil {
		return nil, err
	}
	storage, err := getStorage(fc.StorageType, fc.StorageConfig)
	if err != nil {
		return nil, err
	}
	ff, err := formats.GetFileFormat(fc.Format, storage)
	if err != nil {
		return nil, err
	}
	return &FileConnector{
		Config:     fc,
		Storage:    storage,
		FileFormat: &ff,
	}, nil
}
