package dummy

import (
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

const TIMEOUT = 3 * time.Second

func TestNew(t *testing.T) {
	config := connectors.ConnectorConfig{}
	c := New(config)
	if c == nil {
		t.Error("expected a connector, got nil")
	}
}

func TestRead(t *testing.T) {
	recordCount := 100
	config := connectors.ConnectorConfig{}
	c := New(config)
	p := test.GetExtParams()
	w := wire.New()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	params := map[string]interface{}{
		"ext_params":   p,
		"record_count": recordCount,
	}
	go func(wg *sync.WaitGroup, params map[string]interface{}, w wire.Wire) {
		defer wg.Done()
		test.TestRead(c, w, params, t)
	}(wg, params, w)

	go func(p meta.ExtParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		if err := test.GenerateTestData(recordCount, ch); err != nil {
			t.Error(err)
		}
	}(p)
	wg.Wait()
}

func TestWrite(t *testing.T) {
	recordCount := 2
	config := connectors.ConnectorConfig{}
	c := New(config)
	p := test.GetExtParams()
	w := wire.New()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, w wire.Wire) {
		defer wg.Done()
		writeCount := 0
		for {
			select {
			case feed := <-w.WriteProgressFeeds():
				progress := feed.Content.(feeds.ProgressMessage)
				writeCount += progress.Count()
			case <-w.IsWriteDone():
				if writeCount != recordCount {
					t.Errorf("expected %d records, got %d", recordCount, writeCount)
				}
				return
			case <-time.After(TIMEOUT):
				t.Errorf("write timeout; expected %d records, got %d", recordCount, writeCount)
				return
			}
		}
	}(wg, w)

	go func(c connectors.OutputConnector, p meta.ExtParams, w wire.Wire) {
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(c, p, w)

	go func(c connectors.InputConnector, p meta.ExtParams, w wire.Wire) {
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(c, p, w)

	go func(p meta.ExtParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		if err := test.GenerateTestData(recordCount, ch); err != nil {
			t.Error(err)
		}
	}(p)

	wg.Wait()
}
