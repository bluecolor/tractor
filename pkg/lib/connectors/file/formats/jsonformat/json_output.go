package jsonformat

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
)

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

func (f *JsonFormat) write(filename string, data []msg.Record, d types.Dataset, wi int) (err error) {
	// todo
	return
}

// todo add batch size, buffer
// todo add timeout
func (f *JsonFormat) StartWriteWorker(filename string, d types.Dataset, w *wire.Wire, wi int) (err error) {
	f.storage.Create(filename)
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
func (f *JsonFormat) Write(d types.Dataset, w *wire.Wire) (err error) {
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
