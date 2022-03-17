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
	connectors.Base
	Config     FileConfig
	FileFormat formats.FileFormat
	Storage    types.Storager
}

func New(config connectors.ConnectorConfig) (*FileConnector, error) {
	fc := FileConfig{}
	if err := config.LoadConfig(&fc); err != nil {
		return nil, err
	}
	return &FileConnector{
		Config: fc,
	}, nil
}

func (f *FileConnector) Connect() error {
	storage, err := getStorage(f.Config.StorageType, f.Config.StorageConfig)
	if err != nil {
		return err
	}
	ff, err := formats.GetFileFormat(f.Config.Format, storage)
	if err != nil {
		return err
	}
	f.FileFormat = ff
	f.Storage = storage

	return nil
}
func (f *FileConnector) Close() error {
	return nil
}
func (f *FileConnector) GetPath(filename string) string {
	url := f.Config.StorageConfig.GetURL()
	if url[len(url)-1] != '/' {
		url += "/"
	}
	return url + filename
}

func init() {
	connectors.Add("file", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config)
	})
}
