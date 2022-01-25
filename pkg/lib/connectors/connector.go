package connectors

import (
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/feed"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Connector interface {
	Connect() error
	Close() error
}

type MetaFinder interface {
	Connector
	FindDatasets(pattern string) ([]meta.Dataset, error)
}

type InputConnector interface {
	Connector
	Read(e meta.ExtInput, w wire.Wire) error
}
type OutputConnector interface {
	Connector
	Write(e meta.ExtOutput, w wire.Wire) error
}

type ParallelWriter interface {
	GetParallel() int
	StartWorker(e meta.ExtOutput, w wire.Wire, i int) error
}

func ParallelWrite(p ParallelWriter, e meta.ExtOutput, w wire.Wire) (err error) {
	parallel := p.GetParallel()
	if parallel < 2 {
		err = p.StartWorker(e, w, 0)
		if err != nil {
			w.SendFeed(feed.NewErrorFeed(feed.SenderOutputConnector, err))
		} else {
			w.SendFeed(feed.NewSuccessFeed(feed.SenderOutputConnector))
		}
		return
	}
	var wg sync.WaitGroup
	for i := 1; i <= parallel; i++ {
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			err := p.StartWorker(e, w, i)
			if err != nil {
				w.SendFeed(feed.NewErrorFeed(feed.SenderOutputConnector, err))
			} else {
				w.SendFeed(feed.NewSuccessFeed(feed.SenderOutputConnector))
			}
		}(&wg, i)
		wg.Add(1)
	}
	wg.Wait()
	return nil
}
