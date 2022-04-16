package file

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/rs/zerolog/log"
	tp "go.beyondstorage.io/v5/types"
)

func (c *FileConnector) ListFiles(options map[string]interface{}) ([]string, error) {
	if c.Storage == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}
	files := []string{}
	path := options["path"].(string)
	if path[len(path)-1] != '/' {
		files := append(files, path)
		return files, nil
	}
	it, err := c.Storage.List(options["path"].(string))
	if err != nil {
		return nil, err
	}
	for {
		o, err := it.Next()
		if errors.Is(err, tp.IterateDone) {
			break
		}
		files = append(files, o.Path)
	}
	return files, nil
}

func (c *FileConnector) GetDataset(options map[string]interface{}) (*types.Dataset, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}
	files, err := c.ListFiles(options)
	if err != nil {
		return nil, err
	}
	dataset := &types.Dataset{}
	dataset.Config = types.Config{
		"files": files,
	}
	options["files"] = files
	fields, err := c.FileFormat.ReadFields(options)
	if err != nil {
		log.Error().Err(err).Msg("failed to get fields")
		return nil, err
	}
	for _, f := range fields {
		dataset.Fields = append(dataset.Fields, &types.Field{
			Name: f,
			Type: "string",
		})
	}
	return dataset, nil
}
