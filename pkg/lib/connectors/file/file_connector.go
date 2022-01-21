package file

import (
	"encoding/json"
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/providers"
	"github.com/bluecolor/tractor/pkg/lib/providers/localfs"
	"github.com/bluecolor/tractor/pkg/models"
)

type FileConfig struct {
	ProviderType string            `json:"providerType"`
	Provider     map[string]string `json:"provider"`
}

type FileConnector struct {
	Provider providers.Provider
}

func NewFileConnector(config string) (*FileConnector, error) {
	var p providers.Provider
	var err error
	fileConfig := FileConfig{}
	if err = json.Unmarshal([]byte(config), &fileConfig); err != nil {
		return nil, err
	}
	switch fileConfig.ProviderType {
	case "localfs":
		p, err = localfs.NewLocalFS(fileConfig.Provider)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("provider type not supported yet " + fileConfig.ProviderType)
	}
	return &FileConnector{
		Provider: p,
	}, nil
}

func (f *FileConnector) FetchDatasetsWithPattern(pattern string) ([]models.Dataset, error) {
	return f.Provider.FetchDatasetsWithPattern(pattern)
}
