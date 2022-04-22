package jsonformat

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"go.beyondstorage.io/v5/pairs"
)

func (f *JsonFormat) Work(filename string, d types.Dataset, w *wire.Wire, wi int) (err error) {
	var buf bytes.Buffer
	var contentstr string
	var content []map[string]interface{}
	size, offset := int64(1000), int64(0) // todo size from .env
	rest := []byte{}
	var readBytes int64 = -1
	bw := wire.NewBuffered(w, d.GetBufferSize())

	for {
		if readBytes != 0 {
			readBytes, err = f.storage.Read(filename, &buf, pairs.WithOffset(offset), pairs.WithSize(size))
			if err != nil {
				return err
			}
			offset += readBytes
			contentstr = contentstr + buf.String()
			buf = bytes.Buffer{}
		}
		if readBytes == 0 && len(rest) == 0 {
			break
		}
	}
	if err := json.Unmarshal([]byte(contentstr), &content); err != nil {
		return err
	}
	for _, item := range content {
		record := msg.Record{}
		for _, field := range d.Fields {
			keys := strings.Split(field.Name, ".")
			if len(keys) == 1 {
				record = append(record, item[field.Name])
			} else {
				var value map[string]interface{} = item
				var result interface{}
				for i, key := range keys {
					if i == len(keys)-1 {
						result = value[key]
					} else {
						value = value[key].(map[string]interface{})
					}
				}
				record = append(record, result)
			}
		}
		bw.Send(record)
	}
	bw.Flush()

	return nil
}
func (f *JsonFormat) StartReadWorker(files []string, p types.Dataset, w *wire.Wire, wi int) (err error) {
	for _, file := range files {
		if err = f.Work(file, p, w, wi); err != nil {
			return
		}
	}
	return
}
func (f *JsonFormat) Read(d types.Dataset, w *wire.Wire) (err error) {
	chunks := utils.ToChunks(d.Config.GetStringArray("files", []string{}), d.GetParallel())
	if err != nil {
		return err
	}
	mwg := esync.NewWaitGroup(w, types.InputConnector)
	for i, chunk := range chunks {
		mwg.Add(1)
		go func(wg *esync.WaitGroup, chunk []string, w *wire.Wire, i int) {
			defer wg.Done()
			if err := f.StartReadWorker(chunk, d, w, i); err != nil {
				mwg.HandleError(err)
			}
		}(mwg, chunk, w, i)
	}
	mwg.Wait()
	return
}
