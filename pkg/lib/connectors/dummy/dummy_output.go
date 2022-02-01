package dummy

import (
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getOutputChannel(p meta.ExtParams) chan<- feeds.Data {
	return p.GetInputDataset().Config.GetChannel(OutputChannelKey)
}

func (c *DummyConnector) Write(p meta.ExtParams, w wire.Wire) (err error) {
	return
}
