package localfs

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/rs/zerolog/log"
)

const TIMEOUT = 3 * time.Second

type file struct {
	Name     string
	Conetent string
}

var csvfiles = []file{
	{
		Name: "test1.csv",
		Conetent: `
			"id","name","age"
			"1","John","30"
			"2","Mary","25"
			"3","Mike","20"
			"4","Linda","35"
		`,
	},
	{
		Name: "test2.csv",
		Conetent: `
			"id","name","age"
			"1","Andrew","20"
			"2","Laeddis","32"
			"3","Dolores","21"
			"4","Chanal","20"
		`,
	},
}

func prepareFiles(dir string, files []file) error {
	for _, f := range files {
		content := strings.TrimSpace(utils.Dedent(f.Conetent))
		err := ioutil.WriteFile(dir+"/"+f.Name, []byte(content), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
func findRecordCount(config CsvConfig, files []file) int {
	total, subHeader := 0, 0
	if config.Header {
		subHeader = 1
	}
	for _, file := range files {
		content := strings.TrimSpace(utils.Dedent(file.Conetent))
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Trim(line, " ") != "" {
				total++
			}
		}
		total -= subHeader
	}
	return total
}
func containsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func TestFindDatasets(t *testing.T) {

	dir, err := ioutil.TempDir("/tmp", "tractor")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	provider, err := NewLocalFSProvider(map[string]interface{}{
		"path": dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	files := []string{"abc.csv", "abcD.csv", "abcDe.csv"}
	for _, file := range files {
		err = ioutil.WriteFile(dir+"/"+file, []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		pattern string
		result  []string
	}{
		{"", files},
		{"^abc", files},
		{"^abcD", []string{"abcD.csv", "abcDe.csv"}},
	}

	for _, test := range tests {
		datasets, err := provider.FindDatasets(test.pattern)
		if err != nil {
			t.Fatal(err)
		}
		if len(datasets) != len(test.result) {
			t.Fatalf("Expected %d datasets, got %d", len(test.result), len(datasets))
		}
		names := []string{}
		for _, dataset := range datasets {
			names = append(names, dataset.Name)
		}
		for _, file := range test.result {
			if !containsStr(names, file) {
				t.Fatalf("Expected dataset %s to be in %v", file, names)
			}
		}
	}
}

func TestCsvRead(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "tractor")
	folderName := strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)
	if err := prepareFiles(dir, csvfiles); err != nil {
		t.Error(err)
	}

	provider, err := NewLocalFSProvider(map[string]interface{}{
		"path": "/tmp",
	})
	ei := meta.ExtInput{
		Dataset: meta.Dataset{
			Name: folderName,
			Config: map[string]interface{}{
				"header": true,
				"quoted": true,
			},
		},
		Parallel: 1,
	}
	w := wire.NewWire()
	go func() {
		if err = provider.Read("csv", ei, w); err != nil {
			t.Error(err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		log.Debug().Msg("waiting for feed")
		defer wg.Done()
		for {
			select {
			case <-w.IsReadDone():
				return
			case feed := <-w.FeedChannel:
				log.Debug().Msgf("got feed  %v", feed.Type)
				if feed.Type == feeds.SuccessFeed && feed.Sender == feeds.SenderInputConnector {
					return
				} else if feed.Type == feeds.ErrorFeed {
					t.Error(feed.Content)
				}
			case <-time.After(TIMEOUT):
				t.Error("timeout no success feed received")
			}
		}
	}(wg, w)

	wg.Add(1)
	csvConfig := CsvConfig{}
	if err := utils.MapToStruct(ei.Dataset.Config, &csvConfig); err != nil {
		t.Error(err)
	}
	expectedRecordCount := findRecordCount(csvConfig, csvfiles)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		log.Debug().Msg("waiting for data")
		dataReceived := 0
		for {
			log.Debug().Msgf("data received: %d", dataReceived)
			select {
			case feed := <-w.DataChannel:
				if feed == nil {
					t.Error("no data")
				} else {
					dataReceived += len(feed)
				}
			case <-w.IsReadDone():
				if dataReceived < expectedRecordCount {
					t.Errorf("missing data expected %d got %d", expectedRecordCount, dataReceived)
				} else if dataReceived > expectedRecordCount {
					t.Errorf("too much data expected %d got %d", expectedRecordCount, dataReceived)
				}
				return
			case <-time.After(TIMEOUT):
				t.Error("timeout no success message received")
			}
		}
	}(wg, w)

	wg.Wait()
}
