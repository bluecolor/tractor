package csv

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/feed"
	"github.com/bluecolor/tractor/lib/plugins/outputs"
	"github.com/bluecolor/tractor/lib/utils"
	"github.com/bluecolor/tractor/lib/wire"
)

type Csv struct {
	Path        string          `yaml:"path"`
	File        string          `yaml:"file"`
	ColumnDelim string          `yaml:"column_delim"`
	RecordDelim string          `yaml:"record_delim"`
	Parallel    int             `yaml:"parallel"`
	Header      bool            `yaml:"header"`
	Catalog     *config.Catalog `yaml:"catalog"`
}

var sampleConfig = `
    path: folder to export data
    file: file name, suffix will be added if parallel
    column_delim: column delimiter
    record_delim: record delimiter
	parallel: number of parallel workers
	header: if true will add header to output file
`

func (c *Csv) Description() string {
	return "Write to csv file(s)"
}

func (c *Csv) SampleConfig() string {
	return sampleConfig
}

func (c *Csv) startWorker(w *wire.Wire, i int) error {
	file := func(i int) string {
		switch {
		case i == 0:
			return fmt.Sprintf("%s.csv", c.File)
		default:
			return fmt.Sprintf("%s_%d.csv", c.File, i)
		}
	}(i)
	name := path.Join(c.Path, file)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)

	if c.Header && c.Catalog != nil {
		for _, field := range c.Catalog.Fields {
			writer.WriteString(field.Name)
			writer.WriteString(c.ColumnDelim)
		}
		writer.WriteString(c.RecordDelim)
	}
	for data := range w.ReadData() {
		buff := c.dataToString(data)
		_, err := writer.WriteString(fmt.Sprintf("%s%s", buff, c.RecordDelim))
		if err != nil {
			return err
		}
		progress := feed.NewWriteProgress(len(data))
		w.SendFeed(progress)
	}
	writer.Flush()
	return nil
}

func (c *Csv) Write(w *wire.Wire) (err error) {
	if c.Parallel < 2 {
		return c.startWorker(w, 0)
	}
	var wg sync.WaitGroup

	for i := 1; i <= c.Parallel; i++ {
		go func(wg *sync.WaitGroup) {
			c.startWorker(w, i)
			wg.Done()
		}(&wg)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}

func newCsv(options map[string]interface{}, catalog *config.Catalog) *Csv {
	csv := &Csv{
		ColumnDelim: ",",
		RecordDelim: "\n",
		Parallel:    1,
		File:        "",
		Header:      false,
		Catalog:     catalog,
	}
	utils.ParseOptions(csv, options)
	if csv.File == "" && csv.Catalog != nil && csv.Catalog.Name != "" {
		csv.File = csv.Catalog.Name
	}
	return csv
}

func init() {
	outputs.Add("csv", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (outputs.OutputPlugin, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newCsv(options, catalog), nil
	})
}
