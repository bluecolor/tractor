package csvformat

import (
	"bytes"
	"encoding/csv"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"go.beyondstorage.io/v5/pairs"
)

func getInputCsvDelimiter(p meta.ExtParams) string {
	return p.GetInputDataset().Config.GetString(DelimiterKey, ",")
}
func getLazyQuotes(p meta.ExtParams) bool {
	return p.GetInputDataset().Config.GetBool(QuotesKey, true)
}
func getInputFiles(p meta.ExtParams) []string {
	return p.GetInputDataset().Config.GetStringArray(FilesKey, []string{})
}
func getHeader(p meta.ExtParams) bool {
	return p.GetInputDataset().Config.GetBool(HeaderKey, true)
}

func (f *CsvFormat) Work(filename string, p meta.ExtParams, w wire.Wire, wi int) (err error) {
	var buf bytes.Buffer
	size, offset := int64(1000), int64(0) // todo size from .env
	rest := []byte{}
	records := []feeds.Record{}
	var lines []string
	var isFirstRecord = true
	var hasHeader = getHeader(p)
	var readBytes int64 = -1
	for {
		if readBytes != 0 {
			readBytes, err = f.storage.Read(filename, &buf, pairs.WithOffset(offset), pairs.WithSize(size))
			offset += readBytes
			if err != nil {
				return err
			}
		}
		if readBytes == 0 && len(rest) == 0 {
			break
		}
		if len(rest) > 0 {
			rest = append(rest, []byte("\n")...)
			buf = *bytes.NewBuffer(append(rest, buf.Bytes()...))
		}
		lines, rest = toLinesWithRest(buf.String())
		csvReader := csv.NewReader(strings.NewReader(strings.Join(lines, "\n")))
		csvReader.Comma = []rune(getInputCsvDelimiter(p))[0]
		csvReader.LazyQuotes = getLazyQuotes(p)
		rows, err := csvReader.ReadAll()
		if err != nil {
			return err
		}
		for _, row := range rows {
			if isFirstRecord && hasHeader {
				isFirstRecord = false
				continue
			}
			record, err := toRecord(row, p.GetInputDataset().Fields)
			if err != nil {
				return err
			}
			records = append(records, record)
			if len(records) >= p.GetInputBufferSize() {
				w.SendData(records)
				w.SendFeed(feeds.NewReadProgress(len(records)))
				records = []feeds.Record{}
			}
		}
		buf.Reset()
	}
	if len(records) > 0 {
		w.SendData(records)
		w.SendFeed(feeds.NewReadProgress(len(records)))
	}
	return nil
}
func (f *CsvFormat) StartReadWorker(files []string, p meta.ExtParams, w wire.Wire, wi int) (err error) {
	for _, file := range files {
		if err = f.Work(file, p, w, wi); err != nil {
			return
		}
	}
	return
}
func (f *CsvFormat) Read(p meta.ExtParams, w wire.Wire) (err error) {
	chunks, err := getFileChunks(getInputFiles(p), p.GetInputParallel())
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for i, chunk := range chunks {
		wg.Add(1)
		go func(wg *sync.WaitGroup, chunk []string, w wire.Wire, i int) {
			defer wg.Done()
			if err := f.StartReadWorker(chunk, p, w, i); err != nil {
				w.SendFeed(feeds.NewError(feeds.SenderInputConnector, err))
			}
		}(wg, chunk, w, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccess(feeds.SenderInputConnector))
	w.ReadDone()
	return
}
