package csvformat

import (
	"fmt"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

func generateHeader(fields []*types.Field, delimiter string) string {
	var header []string

	for _, field := range fields {
		header = append(header, field.Name)
	}
	return strings.Join(header, delimiter)
}
func generateFileNames(filename string, parallel int) []string {
	if parallel == 1 {
		return []string{filename}
	}
	name, extention := utils.SplitExt(filename)
	var fileNames []string
	for i := 0; i < parallel; i++ {
		fileNames = append(fileNames, name+"_"+fmt.Sprint(i)+extention)
	}
	return fileNames
}

func (f *CsvFormat) write(filename string, data []msg.Record, d types.Dataset, wi int) (err error) {
	buffer := make([][]string, len(data))
	fields := d.GetFields()

	for i, r := range data {
		for j, _ := range fields {
			colval := r[j]
			buffer[i] = append(buffer[i], colval.(string))
		}
	}
	if len(buffer) == 0 {
		return
	}
	lines := make([]string, len(buffer))
	for i, row := range buffer {
		lines[i] = strings.Join(row, d.Config.GetString("delimiter", ","))
	}
	content := strings.Join(lines, "\n") + "\n"
	_, err = f.storage.Write(filename, strings.NewReader(content), int64(len(content)))
	if err != nil {
		log.Error().Err(err).Msg("failed to write data")
		return err
	}
	return
}

// todo add batch size, buffer
// todo add timeout
func (f *CsvFormat) StartWriteWorker(filename string, d types.Dataset, w *wire.Wire, wi int) (err error) {
	f.storage.Create(filename)
	header := generateHeader(d.GetFields(), d.Config.GetString("delimiter", ",")) + "\n"
	f.storage.Write(filename, strings.NewReader(header), int64(len(header)))
	for {
		data, ok := <-w.ReceiveData()
		if !ok {
			return nil
		}
		if err := f.write(filename, data, d, wi); err != nil {
			return err
		}
		w.SendOutputProgress(data.Count())
	}
}
func (f *CsvFormat) Write(d types.Dataset, w *wire.Wire) (err error) {
	var parallel int = d.GetParallel()
	files := generateFileNames(d.Config.GetString("file", "out.csv"), parallel)
	mwg := esync.NewWaitGroup(w, types.OutputConnector)
	for i, file := range files {
		mwg.Add(1)
		go func(wg *esync.WaitGroup, file string, w *wire.Wire, wi int) {
			defer wg.Done()
			if err := f.StartWriteWorker(file, d, w, wi); err != nil {
				w.SendOutputError(err)
			}
		}(mwg, file, w, i)
	}
	mwg.Wait()
	return
}
