package dummy

import (
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

const TIMEOUT = 30000 * time.Second

func TestNew(t *testing.T) {
	config := connectors.ConnectorConfig{}
	_, err := New(config)
	if err != nil {
		t.Error(err)
	}
}

func TestIO(t *testing.T) {
	t.Parallel()
	recordCount := 10
	config := connectors.ConnectorConfig{}
	connector, err := New(config)
	if err != nil {
		t.Error(err)
	}
	p := test.GetSession()
	w := wire.New()
	wg := &sync.WaitGroup{}

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w)
		memo := casette.Memo()
		if memo.HasError() {
			for _, e := range memo.Errors() {
				t.Error(e)
			}
		}
		if memo.ReadCount() != recordCount {
			t.Errorf("(read) expected %d records, got %d", recordCount, memo.ReadCount())
		}
		if memo.WriteCount() != recordCount {
			t.Errorf("(write) expected %d records, got %d", recordCount, memo.WriteCount())
		}
	}(wg, t)

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, p types.SessionParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestData(recordCount, ch); err != nil {
			t.Error(err)
		}
	}(wg, p)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	// start input connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Input, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	wg.Wait()
}

func TestParallelIO(t *testing.T) {
	t.Parallel()
	recordCount := 100
	config := connectors.ConnectorConfig{}
	connector, err := New(config)
	if err != nil {
		t.Error(err)
	}
	p := test.GetSession().WithInputParallel(2).WithOutputParallel(2)
	w := wire.New()
	wg := &sync.WaitGroup{}

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w)
		memo := casette.Memo()
		if memo.HasError() {
			for _, e := range memo.Errors() {
				t.Error(e)
			}
		}
		if memo.ReadCount() != recordCount {
			t.Errorf("(read) expected %d records, got %d", recordCount, memo.ReadCount())
		}
		if memo.WriteCount() != recordCount {
			t.Errorf("(write) expected %d records, got %d", recordCount, memo.WriteCount())
		}
	}(wg, t)

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, p types.SessionParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestDataWithDuration(recordCount, ch, 3*time.Second); err != nil {
			t.Error(err)
		}
	}(wg, p)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	// start input connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Input, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	wg.Wait()
}

func TestIOError(t *testing.T) {
	t.Parallel()
	recordCount := 100
	config := connectors.ConnectorConfig{}
	connector, err := New(config)
	if err != nil {
		t.Error(err)
	}
	p := test.GetSession().WithInputParallel(2).WithOutputParallel(2)
	w := wire.New()
	wg := &sync.WaitGroup{}

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w)
		memo := casette.Memo()
		if !memo.HasInputError() {
			t.Error("expected an input error")
			return
		}
		if memo.HasInputSuccess() {
			t.Error("unexpected an input success")
			return
		}
	}(wg, t)

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, p types.SessionParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestDataWithDuration(recordCount, ch, 4*time.Second); err != nil {
			t.Error(err)
		}
	}(wg, p)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	// start input connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Input, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	go func() {
		// close data channel, this should cause an error in the input connector
		time.Sleep(2 * time.Second)
		w.CloseData()
	}()

	wg.Wait()
}

func TestDataGenConf(t *testing.T) {
	t.Parallel()
	recordCount := 10

	config := connectors.ConnectorConfig{
		"generateFakeData":   true,
		"fakeRecordCount":    recordCount,
		"fakeRecordInterval": 100,
	}
	connector, err := New(config)
	if err != nil {
		t.Error(err)
	}
	p := test.GetSession()
	w := wire.New()
	wg := &sync.WaitGroup{}

	// collect test results
	wg.Add(1)
	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		casette := test.Record(w)
		memo := casette.Memo()
		if memo.HasError() {
			for _, e := range memo.Errors() {
				t.Error(e)
			}
		}
		if memo.ReadCount() != recordCount {
			t.Errorf("(read) expected %d records, got %d", recordCount, memo.ReadCount())
		}
		if memo.WriteCount() != recordCount {
			t.Errorf("(write) expected %d records, got %d", recordCount, memo.WriteCount())
		}
	}(wg, t)

	// start output connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Output, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Write(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	// start input connector
	wg.Add(1)
	go func(wg *sync.WaitGroup, c connectors.Input, p types.SessionParams, w *wire.Wire) {
		defer wg.Done()
		if err := c.Read(p, w); err != nil {
			t.Error(err)
		}
	}(wg, connector, p, w)

	wg.Wait()
}
