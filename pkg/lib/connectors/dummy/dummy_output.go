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
		select {
		case <-w.Context().Done():
			if err := w.Context().Err(); err != nil {
				w.SendOutputCancelled(err)
				return err
			}
			return nil
		case msg, ok := <-w.GetDataMessage():
			if !ok {
				w.SendOutputSuccess()
				return nil
			}
			od, err := meta.ToOutputData(msg.Data(), p)
			if err != nil {
				w.SendOutputError(err)
				return err
			}
			outputChannel <- od
			w.SendOutputProgress(msg.Count())
		}
	}
}
