package csvformat

import (
	"os"
	"strings"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

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
		fileNames = append(fileNames, name+"_"+string(i)+extention)
	}
	return fileNames
}

func (f *CsvFormat) write(data feeds.Data) (err error) {
	// var line []string
	// for _, field := range data.Fields {
	// 	line = append(line, field.Value)
	// }
	// f.writer.Write([]byte(strings.Join(line, f.csvconfig.Delimiter)))
	return nil
}

func (f *CsvFormat) StartWriteWorker(filename string, p meta.ExtParams, w wire.Wire) (err error) {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// header := generateHeader(, o.csvconfig)

	for data := range w.ReadData() {
		if data == nil {
			break
		}
		err := f.write(data)
		if err != nil {
			return err
		}
	}
	w.WriteWorkerDone()
	return nil
}

func (f *CsvFormat) Write(p meta.ExtParams, w wire.Wire) (err error) {
	var parallel int = p.GetOutputParallel()
	if parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", parallel, 1)
	}
	files := generateFileNames(p.GetOutputDataset().Name, parallel)

	wg := &sync.WaitGroup{}
	for i, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, file string, w wire.Wire, wi int) {
			defer wg.Done()
			err := f.StartWriteWorker(file, p, w)
			if err != nil {
				w.SendFeed(feeds.NewErrorFeed(feeds.SenderOutputConnector, err))
			}
		}(wg, file, w, i)
	}
	wg.Wait()
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderOutputConnector))
	w.WriteDone()
	return
}
