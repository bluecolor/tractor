package localfs

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

type CsvConfig struct {
	Delimiter string `json:"delimiter"`
	Header    bool   `json:"header"`
	Quoted    bool   `json:"quoted"`
}

func (f *LocalFSProvider) ReadCsvWorker(
	filepath string, csvConfig CsvConfig, ec meta.Config, fields []meta.Field, w wire.Wire, i int,
) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	r := csv.NewReader(bufio.NewReader(file))
	delimiter := csvConfig.Delimiter
	if delimiter == "" {
		delimiter = ","
	}
	r.Comma = []rune(delimiter)[0]
	r.LazyQuotes = !csvConfig.Quoted
	bufferSize := ec.GetInt("buffer_size", 100)
	buffer := []feeds.Record{}
	isFirst := true
	for {
		csvRec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if csvConfig.Header && isFirst {
			isFirst = false
			continue
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
	files, ok := e.Dataset.Config["files"].([]string)
	if !ok {
		return errors.New("files not found")
	}
	parallel := e.Parallel
	chunkCount := len(files) / parallel
	chunk := []string{}

	wg := &sync.WaitGroup{}
	for i, j := 0, 0; i < len(files); i++ {
		chunk = append(chunk, files[i])
		if (i%chunkCount == 0 && i != 0) || i == len(files)-1 {
			j++
			wg.Add(1)
			go func(wg *sync.WaitGroup, chunk []string, j int) {
				defer wg.Done()
				for _, file := range chunk {
					log.Debug().Msg(fmt.Sprintf("read csv file %s", file))
					if err := f.ReadCsvWorker(file, csvConfig, e.Config, e.Fields, w, j); err != nil {
						w.SendFeed(feeds.NewErrorFeed(feeds.SenderInputConnector, err))
					}
					log.Debug().Msgf("read csv file %s done", file)
				}
			}(wg, chunk, j)
			chunk = []string{}
		}
	}
	wg.Wait()
	return
}
