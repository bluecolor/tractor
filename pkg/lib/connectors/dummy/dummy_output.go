package dummy

import (
	"context"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getOutputChannel(p params.SessionParams) chan<- interface{} {
	return p.GetOutputDataset().Config.GetChannel(OutputChannelKey)
}

func (c *DummyConnector) StartWriteWorker(ctx context.Context, channel chan<- interface{}, w *wire.Wire) (err error) {
	defer func() {
		if e := recover(); err != nil {
			err = e.(error)
		}
	}()
	for {
		select {
		case data, ok := <-w.ReceiveData():
			if !ok {
				return nil
			}
			channel <- data
			w.SendOutputProgress(data.Count())
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *DummyConnector) Write(p params.SessionParams, w *wire.Wire) error {
	var channel chan<- interface{} = getOutputChannel(p)
	var parallel int = p.GetOutputParallel()
	wg := esync.NewWaitGroup(w, types.OutputConnector)
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *esync.WaitGroup, w *wire.Wire) {
			defer wg.Done()
			if err := c.StartWriteWorker(wg.Context(), channel, w); err != nil {
				wg.HandleError(err)
			}
		}(wg, w)
	}
	return wg.Wait()
}
