package localfs

import (
	"io/ioutil"
	"regexp"

	"github.com/bluecolor/tractor/pkg/lib/connectors/file/providers"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

type LocalFSConfig struct {
	Path string `json:"path"`
}

type LocalFSProvider struct {
	config LocalFSConfig
}

func NewLocalFSProvider(config map[string]interface{}) (*LocalFSProvider, error) {
	localFSConfig := LocalFSConfig{}
	if err := utils.MapToStruct(config, &localFSConfig); err != nil {
		return nil, err
	}
	return &LocalFSProvider{
		config: localFSConfig,
	}, nil
}
func (f *LocalFSProvider) FindDatasets(pattern string) ([]meta.Dataset, error) {
	files, err := ioutil.ReadDir(f.config.Path)
	if err != nil {
		return nil, err
	}
	datasets := []meta.Dataset{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if pattern != "" {
			match, _ := regexp.MatchString(pattern, fileName)
			log.Info().Msgf("match: %v, pattern: %v, fileName: %v", match, pattern, fileName)
			if !match {
				continue
			}
		}
		datasets = append(datasets, meta.Dataset{
			Name: fileName,
		})
	}
	return datasets, nil
}

func init() {
	providers.Add("localfs", func(config map[string]interface{}) (providers.Provider, error) {
		return NewLocalFSProvider(config)
	})
}
