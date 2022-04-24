package file

import (
	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/formats"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/file/formats/all"
	"github.com/bluecolor/tractor/pkg/lib/connectors/file/storage"
	"go.beyondstorage.io/v5/types"
)

type FileConfig struct {
	Provider map[string]interface{} `json:"provider"`
	Format   map[string]interface{} `json:"format"`
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
	storage, err := storage.GetStorage(f.Config.Provider)
	if err != nil {
		return err
	}
	ff, err := formats.GetFileFormat(f.Config.Format["code"].(string), storage)
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
	return f.Config.Provider["code"].(string) + "://" + filename
}

func init() {
	connectors.Add("file", func(config connectors.ConnectorConfig) (connectors.Connector, error) {
		return New(config)
	})
}
