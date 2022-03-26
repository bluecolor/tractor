package runner

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/test"
	"github.com/bluecolor/tractor/pkg/lib/types"
)

func TestRunner(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	recordCount := 6
	connection := &types.Connection{
		ConnectionType: "dummy",
	}

	id, od := test.GetDatasets()
	id.Connection = connection
	od.Connection = connection
	id.Config.SetInt("parallel", 2)
	od.Config.SetInt("parallel", 2)
	session := types.Session{
		Extraction: &types.Extraction{
			SourceDataset: &id,
			TargetDataset: &od,
		},
	}
	runner, err := New(context.Background(), session)
	if err != nil {
		t.Error(err)
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

	if err := runner.Run(); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
