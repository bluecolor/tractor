package file

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

const TIMEOUT = 3 * time.Second

var csvTestFiles = []map[string]string{
	{
		"name": "test1.csv",
		"content": `
			"id","name","age"
			"1","John","30"
			"2","Mary","25"
			"3","Mike","20"
			"4","Linda","35"
		`,
	},
	{
		"name": "test2.csv",
		"content": `
			"id","name","age"
			"1","Andrew","20"
			"2","Laeddis","32"
			"3","Dolores","21"
			"4","Chanal","20"
		`,
	},
}

func getFileNames(files []map[string]string) []string {
	var names []string
	for _, f := range files {
		names = append(names, f["name"])
	}
	return names
}
func getRecordCount(files []map[string]string, header bool) (count int) {
	for _, f := range files {
		content := strings.TrimSpace(utils.Dedent(f["content"]))
		count += strings.Count(content, "\n")
		if !header {
			count++
		}
	}
	return
}
func prepareCsvTestFiles(connector *FileConnector) (err error) {
	for _, f := range csvTestFiles {
		content := strings.TrimSpace(utils.Dedent(f["content"]))
		if err = connector.Storage.Delete(f["name"]); err != nil {
			log.Warn().Err(err).Msg("delete file")
		}
		if _, err = connector.Storage.Write(f["name"], strings.NewReader(content), int64(len(content))); err != nil {
			return
		}
	}
	return
}
func prepareFiles(connector *FileConnector) (err error) {
	prepareCsvTestFiles(connector)
	return
}

func testCsvIO(connector *FileConnector, t *testing.T) {
	outfile := "test_out.csv"
	connector.Storage.Delete(outfile)
	fields := []meta.Field{
		{
			Name: "id",
			Type: "string",
		},
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "age",
			Type: "int",
		},
	}
	p := meta.ExtParams{}
	p.WithInputDataset(meta.Dataset{
		Name:   "test",
		Fields: fields,
		Config: meta.Config{
			"files": getFileNames(csvTestFiles),
		},
	}).WithOutputDataset(meta.Dataset{
		Name:   "test",
		Fields: fields,
		Config: meta.Config{
			"file_name": outfile,
		},
	}).WithFieldMappings([]meta.FieldMapping{
		{
			SourceField: fields[0],
			TargetField: fields[0],
		},
		{
			SourceField: fields[1],
			TargetField: fields[1],
		},
		{
			SourceField: fields[2],
			TargetField: fields[2],
		},
	})
	w, _, _ := wire.NewWithDefaultTimeout()

	go func(w *wire.Wire, p meta.ExtParams) {
		if err := connector.Write(p, w); err != nil {
			t.Error(err)
		}
	}(w, p)

	go func(w *wire.Wire, p meta.ExtParams) {
		if err := connector.Read(p, w); err != nil {
			t.Error(err)
		}
	}(w, p)

	expectedRecordCount := getRecordCount(csvTestFiles, true)
	readCount := 0
	writeCount := 0
	var inputSuccess, outputSuccess bool = false, false

	for {
		select {
		case feed, ok := <-w.GetFeedback():
			if ok {
				readCount += feed.InputProgress()
				writeCount += feed.OutputProgress()
				if feed.IsError() {
					t.Error(feed.Error())
					return
				} else if feed.IsInputSuccess() {
					inputSuccess = true
					w.CloseData()
				} else if feed.IsOutputSuccess() {
					outputSuccess = true
				}
				if inputSuccess && outputSuccess {
					if readCount != expectedRecordCount {
						t.Errorf("(input) expected %d records, got %d", expectedRecordCount, readCount)
					}
					if writeCount != expectedRecordCount {
						t.Errorf("(output) expected %d records, got %d", expectedRecordCount, writeCount)
					}
					return
				}
			} else {
				return
			}
		case <-time.After(TIMEOUT):
			log.Debug().Msg("read count " + strconv.Itoa(readCount))
			log.Debug().Msg("write count " + strconv.Itoa(writeCount))
			log.Debug().Msg("input success " + strconv.FormatBool(inputSuccess))
			log.Debug().Msg("output success " + strconv.FormatBool(outputSuccess))
			t.Error("timeout")
			return
		}
	}

}

func TestNewFileConnector(t *testing.T) {
	configs := []connectors.ConnectorConfig{
		{
			"storageType": "fs",
			"format":      "csv",
			"storageConfig": map[string]interface{}{
				"url": "fs:///tmp/",
			},
		},
	}
	for _, config := range configs {
		if _, err := New(config); err != nil {
			t.Error(err)
		}
	}
}
func TestCsvIO(t *testing.T) {
	configs := []connectors.ConnectorConfig{
		{
			"storageType": "fs",
			"format":      "csv",
			"storageConfig": map[string]interface{}{
				"url": "fs:///tmp/",
			},
		},
	}
	for _, config := range configs {
		connector, err := New(config)
		if err != nil {
			t.Error(err)
		}
		if err := connector.Connect(); err != nil {
			t.Error(err)
		}
		prepareFiles(connector)
		testCsvIO(connector, t)
		if err := connector.Close(); err != nil {
			t.Error(err)
		}
	}
}
