package esync

import (
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/rs/zerolog/log"
)

type WaitGroup struct {
	mu    sync.Mutex
	wg    *sync.WaitGroup
	count int
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		wg: &sync.WaitGroup{},
	}
}

func (g *WaitGroup) Count() int {
	return g.count
}
func (g *WaitGroup) Add(n int) {
	g.mu.Lock()
	g.wg.Add(n)
	g.count += n
	g.mu.Unlock()
	log.Debug().Msgf("count is now %d", g.count)
}
func (g *WaitGroup) Done(ct types.ConnectorType) {
	log.Debug().Msgf("try get lock for %s", ct)
	g.mu.Lock()
	log.Debug().Msgf("got lock")
	g.count--
	g.mu.Unlock()
	log.Debug().Msgf("unlocked %v", ct)
	log.Debug().Msgf("++++++++count is now %d", g.count)
	g.wg.Done()

}
func (g *WaitGroup) Wait() {
	g.wg.Wait()
	log.Debug().Msgf("-------Wait group finished")
}
func (g *WaitGroup) Cancel(ct types.ConnectorType) {
	for i := 0; i < g.count; i++ {
		log.Debug().Msgf("!!!!Cancelling wait group with %d tasks", g.count)
		g.Done(ct)
	}
}
