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

func getOutputFileName(p types.SessionParams) string {
	return p.GetOutputDataset().Config.GetString(FileNameKey, p.GetOutputDataset().Name+".csv")
}
func getOutputCsvDelimiter(p types.SessionParams) string {
	return p.GetOutputDataset().Config.GetString(DelimiterKey, ",")
}
func generateHeader(p types.SessionParams) string {
	var header []string
	for _, field := range p.GetOutputDatasetFields() {
		header = append(header, field.Name)
	}
	return strings.Join(header, getOutputCsvDelimiter(p))
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

func (f *CsvFormat) write(filename string, data []msg.Record, p types.SessionParams, wi int) (err error) {
	od := p.GetOutputDataset()
	buffer := make([][]string, len(data))
	for i, r := range data {
		for _, f := range od.Fields {
			colval, ok := r[p.GetSourceFieldNameByTargetFieldName(f.Name)]
			if !ok {
				log.Debug().Msgf("field %s not found in record %d", f.Name, i)
			}
			buffer[i] = append(buffer[i], colval.(string))
		}
	}
	if len(buffer) == 0 {
		return
	}
	lines := make([]string, len(buffer))
	for i, row := range buffer {
		lines[i] = strings.Join(row, getOutputCsvDelimiter(p))
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
func (f *CsvFormat) StartWriteWorker(filename string, p types.SessionParams, w *wire.Wire, wi int) (err error) {
	f.storage.Create(filename)
	header := generateHeader(p) + "\n"
	f.storage.Write(filename, strings.NewReader(header), int64(len(header)))
	for {
		data, ok := <-w.ReceiveData()
		if !ok {
			return nil
		}
		if err := f.write(filename, data, p, wi); err != nil {
			return err
		}
		w.SendOutputProgress(data.Count())
	}
}
func (f *CsvFormat) Write(p types.SessionParams, w *wire.Wire) (err error) {
	var parallel int = p.GetOutputParallel()
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
	}
	files := generateFileNames(getOutputFileName(p), parallel)

	mwg := esync.NewWaitGroup(w, types.OutputConnector)
	for i, file := range files {
		mwg.Add(1)
		go func(wg *esync.WaitGroup, file string, w *wire.Wire, wi int) {
			defer wg.Done()
			if err := f.StartWriteWorker(file, p, w, wi); err != nil {
				w.SendOutputError(err)
			}
		}(mwg, file, w, i)
	}
	mwg.Wait()
	return
}
