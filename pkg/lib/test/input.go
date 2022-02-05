package test

import (
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

func TestRead(connector connectors.Input, w wire.Wire, params map[string]interface{}, t *testing.T) {

	p, ok := params["ext_params"].(meta.ExtParams)
	if !ok {
		t.Error("ext_params is not a meta.ExtParams")
		return
	}
	recordCount, ok := params["record_count"].(int)
	if !ok {
		t.Error("record_count is not an int")
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		log.Debug().Msg("waiting for data")
		dataReceived := 0
		for {
			select {
			case feed := <-w.ReadFeeds():
				log.Debug().Msgf("received data feed %v", feed)
			case feed := <-w.ReadData():
				if feed == nil {
					if dataReceived < recordCount {
						t.Errorf("missing data before timeout expected %d, got %d", recordCount, dataReceived)
					} else if dataReceived > recordCount {
						t.Errorf("too much data before timeout expected %d, got %d", recordCount, dataReceived)
					}
					return
				} else {
					dataReceived += len(feed)
				}
			case <-time.After(TIMEOUT):
				t.Error("timeout before read done")
			}
		}
	}(wg, w)

	go func(p meta.ExtParams, w wire.Wire) {
		if err := connector.Read(p, w); err != nil {
			t.Error(err)
		}
	}(p, w)

	log.Debug().Msg("waiting for done")

	wg.Wait()
}
