package esync

import "sync"

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
	defer g.mu.Unlock()
	g.wg.Add(n)
	g.count += n
}
func (g *WaitGroup) Done() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.wg.Done()
	g.count--
}
func (g *WaitGroup) Wait() {
	g.wg.Wait()
}
func (g *WaitGroup) Cancel() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := 0; i < g.count; i++ {
		g.Done()
	}
}
