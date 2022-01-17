package csv

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/utils"
)

type Csv struct {
	Path        string `yaml:"path"`
	File        string `yaml:"file"`
	ColumnDelim string `yaml:"column_delim"`
	Parallel    int    `yaml:"parallel"`
	Header      bool   `yaml:"header"`
}

var files [][]string

var sampleConfig = `
    path: folder to export data
    file: name of the file to be ingested, leave empty for multifile
    column_delim: column delimiter
    parallel: applies only to multifile
`

func (c *Csv) Description() string {
	return "Read from csv file(s)"
}

func (c *Csv) SampleConfig() string {
	return sampleConfig
}

func (c *Csv) startWorker(wire tractor.Wire, files []string) error {
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			println(err.Error())
			return err
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		scanner := bufio.NewScanner(f)
		var data tractor.Data
		header := false
		for scanner.Scan() {
			if c.Header && !header {
				header = true
				continue
			}
			r := strings.Split(scanner.Text(), c.ColumnDelim)
			var record = make([]interface{}, len(r))
			for i, c := range r {
				record[i] = c
			}
			data = append(data, record)
			if len(data) > 100 {
				wire.SendFeed(tractor.NewReadProgress(len(data)))
				wire.SendData(data)
				data = nil
			}
		}
		if len(data) > 0 {
			wire.SendData(data)
			data = nil
		}
		wire.SendFeed(tractor.NewSuccessFeed(tractor.InputPlugin))
		return err
	}
	return nil
}

func (c *Csv) Read(wire tractor.Wire) error {
	if len(files) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	for i := 0; i < len(files); i++ {
		go func(wg *sync.WaitGroup, i int) {
			c.startWorker(wire, files[i])
			wg.Done()
		}(&wg, i)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}

func (c *Csv) Init() error {
	fls, err := c.getFiles()
	if err != nil {
		return err
	}
	files = make([][]string, c.Parallel)
	cnt := int(len(fls) / c.Parallel)
	for i := 0; i < c.Parallel; i++ {
		pack := make([]string, cnt)
		for j := 0; j < cnt; j++ {
			pack[j] = fls[j*(i+1)]
		}
		files[i] = pack
	}
	var x = 0
	for i := cnt * c.Parallel; i < len(fls); i++ {
		files[x] = append(files[x], fls[i])
		x++
	}
	return nil
}

func newCsv(options map[string]interface{}) *Csv {
	csv := &Csv{
		ColumnDelim: ",",
		Parallel:    1,
		Header:      false,
	}
	utils.ParseOptions(csv, options)
	return csv
}

func init() {
	inputs.Add("csv", func(
		config map[string]interface{},
		catalog *config.Catalog,
		params map[string]interface{},
	) (tractor.Input, error) {
		options, err := utils.MergeOptions(config, params)
		if err != nil {
			return nil, err
		}
		return newCsv(options), nil
	})
}