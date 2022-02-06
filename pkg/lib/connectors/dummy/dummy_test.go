package dummy

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

const TIMEOUT = 3 * time.Second

func TestNew(t *testing.T) {
	config := connectors.ConnectorConfig{}
	connector := New(config)
	if connector == nil {
		t.Error("expected a connector, got nil")
	}
}

func Test(t *testing.T) {
	recordCount := 100
	config := connectors.ConnectorConfig{}
	connector := New(config)
	p := test.GetExtParams()
	w, _, cancel := wire.NewWithTimeout(TIMEOUT)
	wg := &sync.WaitGroup{}

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, c context.CancelFunc, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w, c)
		memo := casette.GetMemo()
		if memo.HasError() {
			for _, e := range memo.Errors {
				t.Error(e)
			}
		}
		if memo.ReadCount != recordCount {
			t.Errorf("(read) expected %d records, got %d", recordCount, memo.ReadCount)
		}
		if memo.WriteCount != recordCount {
			t.Errorf("(write) expected %d records, got %d", recordCount, memo.WriteCount)
		}
	}(wg, cancel, t)

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, p meta.ExtParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestData(recordCount, ch); err != nil {
			t.Error(err)
		}
	}(wg, p)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p meta.ExtParams, w wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	// start input connector
	go func(wg *sync.WaitGroup, c connectors.Input, p meta.ExtParams, w wire.Wire) {
		wg.Add(1)
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	wg.Wait()
}
