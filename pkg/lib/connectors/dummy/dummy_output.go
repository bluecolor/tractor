package dummy

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getOutputChannel(p meta.ExtParams) chan<- interface{} {
	return p.GetOutputDataset().Config.GetChannel(OutputChannelKey)
}

func (c *DummyConnector) Write(p meta.ExtParams, w *wire.Wire) error {
	var outputChannel chan<- interface{} = getOutputChannel(p)
	for {
		data, ok := <-w.GetData()
		if !ok {
			w.SendOutputSuccess()
			return nil
		}
		outputChannel <- data
		w.SendOutputProgress(len(data))
	}
}
