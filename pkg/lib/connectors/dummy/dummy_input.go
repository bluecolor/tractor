package dummy

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/esync"
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getInputChannel(p params.SessionParams) <-chan interface{} {
	return p.GetInputDataset().Config.GetChannel(InputChannelKey)
}

func (c *DummyConnector) StartReadWorker(channel <-chan interface{}, w *wire.Wire) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for {
		data, ok := <-channel
		if !ok {
			return
		}
		w.SendData(data)
	}
}

func (c *DummyConnector) Read(p params.SessionParams, w *wire.Wire) error {
	var parallel int = p.GetInputParallel()
	var channel <-chan interface{}
	if c.config.GenerateFakeData {
		channel = c.Generate()
	} else {
		channel = getInputChannel(p)
		if channel == nil {
			return errors.New("no input channel")
		}
	}
	wg := esync.NewWaitGroup(w, types.InputConnector)
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func(wg *esync.WaitGroup, w *wire.Wire) {
			defer wg.Done()
			if err := c.StartReadWorker(channel, w); err != nil {
				wg.HandleError(err)
			}
		}(wg, w)
	}
	return wg.Wait()
}
