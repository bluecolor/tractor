package runner

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/types"
)

func TestNewRunner(t *testing.T) {
	c := types.Connection{
		ConnectionType: "dummy",
	}
	if _, err := New(context.Background(), &c, &c); err != nil {
		t.Error(err)
	}
}

func TestRunner(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	recordCount := 6
	connection := types.Connection{
		ConnectionType: "dummy",
	}
	runner, err := New(context.Background(), &connection, &connection)
	if err != nil {
		t.Error(err)
	}

	id, od := test.GetDatasets()
	id.Config.SetInt("parallel", 2)
	od.Config.SetInt("parallel", 2)
	extraction := types.Extraction{
		SourceDataset: &id,
		TargetDataset: &od,
	}

	// generate test data
	wg.Add(1)
	go func(wg *sync.WaitGroup, d types.Dataset) {
		ch := d.Config.GetChannel("channel")
		defer close(ch)
		defer wg.Done()
		if err := test.GenerateTestDataWithDuration(recordCount, ch, 3*time.Second); err != nil {
			t.Error(err)
		}
	}(wg, id)

	if err := runner.Run(extraction); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
