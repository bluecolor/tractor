package csv

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/outputs"
	"github.com/bluecolor/tractor/utils"
)

type Csv struct {
	Path        string `yaml:"path"`
	File        string `yaml:"file"`
	ColumnDelim string `yaml:"column_delim"`
	RecordDelim string `yaml:"record_delim"`
	Parallel    int    `yaml:"parallel"`
}

var sampleConfig = `
    path: folder to export data
    file: file name, suffix will be added if parallel
    column_delim: column delimiter
    record_delim: record delimiter
`

func (c *Csv) Description() string {
	return "Write to csv file(s)"
}

func (c *Csv) SampleConfig() string {
	return sampleConfig
}

func (c *Csv) startWorker(wire tractor.Wire, i int) error {
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
	w := bufio.NewWriter(f)
	for data := range wire.ReadData() {
		buff := c.dataToString(data)
		_, err := w.WriteString(fmt.Sprintf("%s%s", buff, c.RecordDelim))
		if err != nil {
			return err
		}
		progress := tractor.NewWriteProgress(len(data))
		wire.SendFeed(progress)
	}
	return nil
}

func (c *Csv) Write(wire tractor.Wire) (err error) {
	if c.Parallel < 2 {
		return c.startWorker(wire, 0)
	}
	var wg sync.WaitGroup

	for i := 1; i <= c.Parallel; i++ {
		go func(wg *sync.WaitGroup) {
			c.startWorker(wire, i)
			wg.Done()
		}(&wg)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}

func (c *Csv) Init() error {
	return nil
}

func newCsv(options map[string]interface{}) *Csv {
	csv := &Csv{
		ColumnDelim: ",",
		RecordDelim: "\n",
		Parallel:    1,
		File:        "out",
	}
	utils.ParseOptions(csv, options)
	return csv
}

func init() {
	outputs.Add("csv", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Output, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newCsv(options), nil
	})
}
