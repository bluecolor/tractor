package file

import (
	"strings"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/test"
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
func testCsvRead(connector *FileConnector, t *testing.T) {
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
	w := wire.NewWire()
	params := map[string]interface{}{
		"ext_params":   p,
		"record_count": getRecordCount(csvTestFiles, true),
	}
	test.TestRead(connector, w, params, t)
}
func testCsvWrite(connector *FileConnector, t *testing.T) {
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
	w := wire.NewWire()

	go func(w wire.Wire, p meta.ExtParams) {
		if err := connector.Write(p, w); err != nil {
			t.Error(err)
		}
	}(w, p)

	go func(w wire.Wire, p meta.ExtParams) {
		if err := connector.Read(p, w); err != nil {
			t.Error(err)
		}
	}(w, p)

	expectedRecordCount := getRecordCount(csvTestFiles, true)
	recordCount := 0

	for {
		select {
		case feed := <-w.WriteProgressFeeds():
			progress := feed.Content.(feeds.ProgressMessage)
			recordCount += progress.Count()
		case <-w.IsWriteDone():
			if recordCount != expectedRecordCount {
				t.Errorf("expected %d records, got %d", expectedRecordCount, recordCount)
			}
			return
		case <-time.After(TIMEOUT):
			t.Error("write timeout")
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

func TestCsvRead(t *testing.T) {
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
		testCsvRead(connector, t)
		if err := connector.Close(); err != nil {
			t.Error(err)
		}
	}
}

func TestCsvWrite(t *testing.T) {
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
		testCsvWrite(connector, t)
		if err := connector.Close(); err != nil {
			t.Error(err)
		}
	}
}
