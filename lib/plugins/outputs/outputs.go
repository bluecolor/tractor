package outputs

import (
	"sync"

	"github.com/bluecolor/tractor/lib/feed"
	"github.com/bluecolor/tractor/lib/plugins"
	"github.com/bluecolor/tractor/lib/wire"
)

type OutputPlugin interface {
	plugins.PluginDescriber
	Write(w *wire.Wire) error
}

type ParallelWriter interface {
	GetParallel() int
	StartWorker(w *wire.Wire, i int) error
}

func ParallelWrite(p ParallelWriter, w *wire.Wire) (err error) {
	parallel := p.GetParallel()
	if parallel < 2 {
		err = p.StartWorker(w, 0)
		if err != nil {
			w.SendFeed(feed.NewErrorFeed(feed.SenderOutputPlugin, err))
		} else {
			w.SendFeed(feed.NewSuccessFeed(feed.SenderOutputPlugin))
		}
		return
	}
	var wg sync.WaitGroup
	for i := 1; i <= parallel; i++ {
		go func(wg *sync.WaitGroup) {
			err := p.StartWorker(w, i)
			if err != nil {
				w.SendFeed(feed.NewErrorFeed(feed.SenderOutputPlugin, err))
			} else {
				w.SendFeed(feed.NewSuccessFeed(feed.SenderOutputPlugin))
			}
			wg.Done()
		}(&wg)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}
