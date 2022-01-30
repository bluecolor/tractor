package csvformat

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"go.beyondstorage.io/v5/pairs"
)

type inputOptions struct {
	csvconfig  csvconfig
	bufferSize int
	fields     []meta.Field
	parallel   int
}

func getInputOptions(e meta.ExtInput) (*inputOptions, error) {
	o := &inputOptions{}
	csvconfig := csvconfig{}
	if err := utils.MapToStruct(e.Dataset.Config, &csvconfig); err != nil {
		return nil, err
	}
	o.csvconfig = csvconfig
	o.bufferSize = e.Config.GetInt("buffer_size", 100)
	o.fields = e.Fields
	o.parallel = e.Parallel
	return o, nil
}

func (f *CsvFormat) Work(filename string, o *inputOptions, w wire.Wire, wi int) error {
	var buf bytes.Buffer
	size, offset := int64(1000), int64(0) // todo size from .env
	rest := []byte{}
	records := []feeds.Record{}
	var lines []string
	for {
		n, err := f.storage.Read(filename, &buf, pairs.WithOffset(offset), pairs.WithSize(size))
		if err != nil {
			return err
		} else if n == 0 {
			break
		}
		if len(rest) > 0 {
			buf = *bytes.NewBuffer(append(rest, buf.Bytes()...))
		}
		lines, rest = toLinesWithRest(buf.String())
		scanner := bufio.NewScanner(strings.NewReader(strings.Join(lines, "\n")))
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				continue
			}
			record, err := lineToRecord(line, o.csvconfig.Delimiter, o.fields)
			if err != nil {
				return err
			}
			records = append(records, record)
			if len(records) >= o.bufferSize {
				w.SendData(records)
				w.SendFeed(feeds.NewReadProgress(len(records)))
				records = []feeds.Record{}
			}
		}
	}
	if len(records) > 0 {
		w.SendData(records)
		w.SendFeed(feeds.NewReadProgress(len(records)))
	}
	return nil
}
func (f *CsvFormat) StartReadWorker(files []string, o *inputOptions, w wire.Wire, wi int) (err error) {
	for _, file := range files {
		if err = f.Work(file, o, w, wi); err != nil {
			return
		}
	}
	return
}
func (f *CsvFormat) Read(e meta.ExtInput, w wire.Wire) (err error) {
	options, err := getInputOptions(e)
	if err != nil {
		return err
	}
	chunks, err := getFileChunks(e.Dataset.Config["files"], options.parallel)
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for i, chunk := range chunks {
		wg.Add(1)
		go func(wg *sync.WaitGroup, chunk []string, w wire.Wire, i int) {
			defer wg.Done()
			if err := f.StartReadWorker(chunk, options, w, i); err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderInputConnector, err))
			}
		}(wg, chunk, w, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderInputConnector))
	w.ReadDone()
	return
}
