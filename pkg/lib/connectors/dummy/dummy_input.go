package dummy

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getInputChannel(p meta.ExtParams) <-chan interface{} {
	return p.GetInputDataset().Config.GetChannel(InputChannelKey)
}

func (c *DummyConnector) Read(p meta.ExtParams, w wire.Wire) (err error) {
	var channel <-chan interface{} = getInputChannel(p)

	for {
		select {
		case <-w.Context().Done():
			err = w.Context().Err()
			if err != nil {
				w.SendInputCancelled(err)
			}
			return
		case data, ok := <-channel:
			if !ok {
				w.SendInputSuccess()
				return
			}
			w.SendData(data)
		}
	}
}
