package csv

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/bluecolor/tractor"
)

type Csv struct {
	Path        string `yaml:"path"`
	File        string `yaml:"file"`
	ColumnDelim string `yaml:"column_delim"`
	Parallel    int    `yaml:"parallel"`
}

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
			return err
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		scanner := bufio.NewScanner(f)
		var data tractor.Data

		for scanner.Scan() {
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
