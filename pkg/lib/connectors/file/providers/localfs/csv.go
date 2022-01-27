package localfs

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
)

type CsvConfig struct {
	Delimiter string `json:"delimiter"`
}

func (f *LocalFSProvider) ReadCsvWorker(filename string, csvConfig CsvConfig, fields []meta.Field, w wire.Wire) (err error) {

	return
}

func (f *LocalFSProvider) ReadCsv(e meta.ExtInput, w wire.Wire) (err error) {
	csvConfig := CsvConfig{}
	if err = utils.MapToStruct(e.Dataset.Config, &csvConfig); err != nil {
		return err
	}
	path := f.config.Path + "/" + e.Dataset.Name
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	r := csv.NewReader(bufio.NewReader(file))
	r.Comma = []rune(csvConfig.Delimiter)[0]
	bufferSize := e.Config.GetInt("buffer_size", 100)
	buffer := []feeds.Record{}
	for {
		csvRec, err := r.Read()
		if err != nil {
			break
		}
		record := feeds.Record{}
		for i, f := range e.Fields {
			if f.Name == "" {
				f.Name = fmt.Sprintf("col%d", i)
			}
			record[f.Name] = csvRec[i]
		}
		if len(buffer) >= bufferSize {
			w.SendData(buffer)
			w.SendFeed(feeds.NewReadProgress(len(buffer)))
			buffer = []feeds.Record{}
		} else {
			buffer = append(buffer, record)
		}
	}
	if len(buffer) > 0 {
		w.SendData(buffer)
		w.SendFeed(feeds.NewReadProgress(len(buffer)))
	}

	return
}
