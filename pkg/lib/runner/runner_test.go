package runner

import (
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/test"
)

func TestNewRunner(t *testing.T) {
	c := meta.Connection{
		ConnectionType: "dummy",
	}
	if _, err := New(c, c); err != nil {
		t.Error(err)
	}
}

func TestRunner(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	recordCount := 6
	connection := meta.Connection{
		ConnectionType: "dummy",
	}
	runner, err := New(connection, connection)
	if err != nil {
		t.Error(err)
	}

	p := test.GetExtParams().WithInputParallel(2).WithOutputParallel(2)

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, p meta.ExtParams) {
		ch := p.GetInputDataset().Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestDataWithDuration(recordCount, ch, 3*time.Second); err != nil {
			t.Error(err)
		}
	}(wg, p)

	if err := runner.Run(p); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
