package dummy

import (
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func getOutputChannel(p meta.ExtParams) chan<- feeds.Data {
	return p.GetOutputDataset().Config.GetChannel("channel")
}

func (c *DummyConnector) Write(p meta.ExtParams, w wire.Wire) (err error) {
	ch := getOutputChannel(p)
	for {
		d := <-w.ReadData()
		if d == nil {
			break
		}
		od, err := meta.ToOutputData(d, p)
		if err != nil {
			w.SendWriteErrorFeed(err)
			return err
		}
		ch <- od
		w.SendWriteProgress(len(od))
	}
	w.WriteDone()
	w.SendOutputSuccessFeed()
	return
}
