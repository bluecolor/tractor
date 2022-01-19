package json

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"encoding/json"

	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/feed"
	"github.com/bluecolor/tractor/lib/plugins/outputs"
	"github.com/bluecolor/tractor/lib/utils"
	"github.com/bluecolor/tractor/lib/wire"
)

type Json struct {
	Parallel int             `yaml:"parallel"`
	Path     string          `yaml:"path"`
	File     string          `yaml:"file"`
	Catalog  *config.Catalog `yaml:"catalog"`
}

func newJson(options map[string]interface{}, sourceCatalog *config.Catalog) *Json {
	json := &Json{
		Parallel: 1,
		File:     "",
	}
	utils.ParseOptions(json, options)
	if json.File == "" && json.Catalog != nil && json.Catalog.Name != "" {
		json.File = json.Catalog.Name
	}
	json.mergeSourceCatalog(sourceCatalog)
	return json
}

func (j *Json) dataToString(data feed.Data) (string, error) {
	buff := make([]string, len(data))
	for i, record := range data {
		err := j.mapRecord(&record)
		if err != nil {
			return "", err
		}
		jsonStr, err := json.Marshal(record)
		if err != nil {
			return "", err
		}
		buff[i] = string(jsonStr)
	}
	return strings.Join(buff, ",\n"), nil
}

func (j *Json) GetParallel() int {
	return j.Parallel
}
func (j *Json) Description() string {
	return "write to json file(s)"
}
func (j *Json) SampleConfig() string {
	return `
		path: folder to export data
		file: file name, suffix will be added if parallel
		parallel: number of parallel workers
	`
}
func (j *Json) StartWorker(w *wire.Wire, i int) error {
	file := func(i int) string {
		switch {
		case i == 0:
			return j.File
		default:
			return fmt.Sprintf("%d_%s", i, j.File)
		}
	}(i)
	name := path.Join(j.Path, file)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString("[")
	if err != nil {
		return err
	}
	isFirstBatch := true
	for data := range w.ReadData() {
		buff, err := j.dataToString(data)
		if err != nil {
			return err
		}
		if isFirstBatch {
			isFirstBatch = false
			_, err = writer.WriteString(fmt.Sprintf("%s%s", "", buff))
		} else {
			_, err = writer.WriteString(fmt.Sprintf("%s%s", ",\n", buff))
		}
		if err != nil {
			return err
		}
		progress := feed.NewWriteProgress(len(data))
		w.SendFeed(progress)
	}
	_, err = writer.WriteString("]")
	if err != nil {
		return err
	}
	err = writer.Flush()
	return err
}
func (j *Json) Write(w *wire.Wire) (err error) {
	return outputs.ParallelWrite(j, w)
}

func init() {
	outputs.Add("json", func(
		config map[string]interface{},
		sourceCatalog *config.Catalog,
		params map[string]interface{},
	) (outputs.OutputPlugin, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newJson(options, sourceCatalog), nil
	})
}
