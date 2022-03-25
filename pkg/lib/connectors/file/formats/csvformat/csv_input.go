package csvformat

import (
	"bytes"
	"encoding/csv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"go.beyondstorage.io/v5/pairs"
)

func (f *CsvFormat) Work(filename string, d types.Dataset, w *wire.Wire, wi int) (err error) {
	var buf bytes.Buffer
	size, offset := int64(1000), int64(0) // todo size from .env
	rest := []byte{}
	var lines []string
	var isFirstRecord = true
	var hasHeader = d.Config.GetBool("header", true)
	var readBytes int64 = -1
	bw := wire.NewBuffered(w, d.GetBufferSize())
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
		csvReader.Comma = []rune(d.Config.GetString("delimiter", ","))[0]
		csvReader.LazyQuotes = d.Config.GetBool("quotes", true)
		rows, err := csvReader.ReadAll()
		if err != nil {
			return err
		}
		for _, row := range rows {
			if isFirstRecord && hasHeader {
				isFirstRecord = false
				continue
			}
			bw.Send(row)
		}
		buf.Reset()
	}
	bw.Flush()
	return nil
}
func (f *CsvFormat) StartReadWorker(files []string, p types.Dataset, w *wire.Wire, wi int) (err error) {
	for _, file := range files {
		if err = f.Work(file, p, w, wi); err != nil {
			return
		}
	}
	return
}
func (f *CsvFormat) Read(d types.Dataset, w *wire.Wire) (err error) {
	chunks, err := getFileChunks(d.Config.GetStringArray("files", []string{}), d.GetParallel())
	if err != nil {
		return err
	}
	mwg := esync.NewWaitGroup(w, types.InputConnector)
	for i, chunk := range chunks {
		mwg.Add(1)
		go func(wg *esync.WaitGroup, chunk []string, w *wire.Wire, i int) {
			defer wg.Done()
			if err := f.StartReadWorker(chunk, d, w, i); err != nil {
				w.SendInputError(err)
			}
		}(mwg, chunk, w, i)
	}
	mwg.Wait()
	return
}
