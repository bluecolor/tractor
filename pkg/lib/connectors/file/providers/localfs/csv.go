package localfs

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
)

type CsvConfig struct {
	Delimiter string `json:"delimiter"`
}

func (f *LocalFSProvider) ReadCsvWorker(
	filepath string, csvConfig CsvConfig, ec meta.Config, fields []meta.Field, w wire.Wire, i int,
) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	r := csv.NewReader(bufio.NewReader(file))
	r.Comma = []rune(csvConfig.Delimiter)[0]
	bufferSize := ec.GetInt("buffer_size", 100)
	buffer := []feeds.Record{}
	for {
		csvRec, err := r.Read()
		if err != nil {
			break
		}
		record := feeds.Record{}
		for i, f := range fields {
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

func (f *LocalFSProvider) ReadCsv(e meta.ExtInput, w wire.Wire) (err error) {
	csvConfig := CsvConfig{}
	if err = utils.MapToStruct(e.Dataset.Config, &csvConfig); err != nil {
		return err
	}
	path := f.config.Path + "/" + e.Dataset.Name
	datasets, err := f.FindDatasets(path)
	if err != nil {
		return err
	}
	parallel := e.Parallel
	chunkCount := len(datasets) / parallel
	chunk := []meta.Dataset{}

	wg := &sync.WaitGroup{}
	for i, j := 0, 0; i < len(datasets); i++ {
		chunk = append(chunk, datasets[i])
		if (i%chunkCount == 0 && i != 0) || i == len(datasets)-1 {
			j++
			wg.Add(1)
			go func(chunk []meta.Dataset, j int) {
				defer wg.Done()
				for _, d := range chunk {
					path := f.config.Path + "/" + d.Name
					if err := f.ReadCsvWorker(path, csvConfig, e.Config, e.Fields, w, j); err != nil {
						w.SendFeed(feeds.NewErrorFeed(feeds.SenderInputConnector, err))
					}
				}
			}(chunk, j)
			chunk = []meta.Dataset{}
		}
	}

	wg.Wait()
	w.SendInputSuccessFeed()
	w.ReadDone()
	return
}
