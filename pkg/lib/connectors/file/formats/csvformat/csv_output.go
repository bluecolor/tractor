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

type outputOptions struct {
	csvconfig     csvconfig
	bufferSize    int
	parallel      int
	fieldMappings []meta.FieldMapping
	
}
type 

func getOutputOptions(e meta.ExtOutput) (*outputOptions, error) {
	o := &outputOptions{}
	csvconfig := csvconfig{}
	if err := utils.MapToStruct(e.Dataset.Config, &csvconfig); err != nil {
		return nil, err
	}
	o.csvconfig = csvconfig
	o.bufferSize = e.Config.GetInt("buffer_size", 100)
	o.parallel = e.Parallel
	return o, nil
}
func generateHeader(fields []meta.Field, csvconfig csvconfig) string {
	var header []string
	for _, field := range fields {
		header = append(header, field.Name)
	}
	return strings.Join(header, csvconfig.Delimiter)
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
	var line []string
	for _, field := range data.Fields {
		line = append(line, field.Value)
	}
	f.writer.Write([]byte(strings.Join(line, f.csvconfig.Delimiter)))
	return nil
}

func (f *CsvFormat) StartWriteWorker(filename string, o *outputOptions, w wire.Wire) (err error) {

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

func (f *CsvFormat) Write(e meta.ExtOutput, w wire.Wire) (err error) {
	var parallel int = 1
	if e.Parallel < 1 {
		log.Warn().Msgf("invalid parallel write setting %d. Using %d", e.Parallel, parallel)
	} else {
		parallel = e.Parallel
	}
	files := generateFileNames(e.Dataset.Name, parallel)
	options, err := getOutputOptions(e)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for i, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, file string, w wire.Wire, wi int) {
			defer wg.Done()
			err := f.StartWriteWorker(file, options, w)
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
