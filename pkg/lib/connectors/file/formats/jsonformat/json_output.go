package jsonformat

import (
	"encoding/json"
	"fmt"
	"strings"

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

func (f *JsonFormat) write(filename string, data []msg.Record, d types.Dataset, wi int, output *[]map[string]interface{}) (err error) {
	// todo convert to correct types
	for _, r := range data {
		record := map[string]interface{}{}
		for i, field := range d.GetFields() {
			keys := strings.Split(field.Name, ".")
			if len(keys) == 1 {
				record[field.Name] = string(r[i].([]byte))
			} else {
				tmp := record
				for k, key := range keys {
					if _, ok := record[key]; !ok {
						tmp[key] = map[string]interface{}{}
					}
					if k == len(keys)-1 {
						tmp[key] = string(r[i].([]byte))
					} else {
						tmp = tmp[key].(map[string]interface{})
					}
				}
			}
		}
		*output = append(*output, record)
	}
	return
}
func (f *JsonFormat) persist(filename string, output []map[string]interface{}) error {
	f.storage.Create(filename)
	ob, err := json.Marshal(output)
	if err != nil {
		return err
	}
	outstr := string(ob)
	_, err = f.storage.Write(filename, strings.NewReader(string(outstr)), int64(len(outstr)))
	if err != nil {
		return err
	}
	return nil
}

// todo add batch size, buffer
// todo add timeout
func (f *JsonFormat) StartWriteWorker(filename string, d types.Dataset, w *wire.Wire, wi int) (err error) {
	output := []map[string]interface{}{}
	for {
		data, ok := <-w.ReceiveData()
		if !ok {
			return f.persist(filename, output)
		}
		if err := f.write(filename, data, d, wi, &output); err != nil {
			return err
		}
		w.SendOutputProgress(data.Count())
	}
}
func (f *JsonFormat) Write(d types.Dataset, w *wire.Wire) (err error) {
	var parallel int = d.GetParallel()
	files := generateFileNames(d.Config.GetString("path", "out.json"), parallel)
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
