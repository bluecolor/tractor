package csvformat

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

func getOutputFileName(p meta.ExtParams) string {
	return p.GetOutputDataset().Config.GetString(FileNameKey, p.GetOutputDataset().Name+".csv")
}
func getOutputCsvDelimiter(p meta.ExtParams) string {
	return p.GetOutputDataset().Config.GetString(DelimiterKey, ",")
}
func generateHeader(p meta.ExtParams) string {
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

func (f *CsvFormat) write(filename string, data feeds.Data, p meta.ExtParams, wi int) (err error) {
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
func (f *CsvFormat) StartWriteWorker(filename string, p meta.ExtParams, w wire.Wire, wi int) (err error) {
	f.storage.Create(filename)
	header := generateHeader(p) + "\n"
	f.storage.Write(filename, strings.NewReader(header), int64(len(header)))

	for {
		data, ok := <-w.ReadData()
		if !ok {
			break
		}
		if err := f.write(filename, data, p, wi); err != nil {
			w.SendWriteErrorFeed(err)
			return err
		}
		w.SendWriteProgress(len(data))
	}
	w.WriteWorkerDone()
	return
}
func (f *CsvFormat) Write(p meta.ExtParams, w wire.Wire) (err error) {
	var parallel int = p.GetOutputParallel()
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
	}
	files := generateFileNames(getOutputFileName(p), parallel)

	wg := &sync.WaitGroup{}
	for i, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, file string, w wire.Wire, wi int) {
			defer wg.Done()
			err := f.StartWriteWorker(file, p, w, wi)
			if err != nil {
				w.SendFeed(feeds.NewError(feeds.SenderOutputConnector, err))
			}
		}(wg, file, w, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccess(feeds.SenderOutputConnector))
	w.WriteDone()
	return
}
