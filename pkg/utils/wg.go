package utils

import (
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type WaitGroup struct {
	mu    sync.Mutex
	wg    *sync.WaitGroup
	count int
	err   error
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
func (g *WaitGroup) CancelWithError(err error) {
	g.err = err
	g.Cancel()
}
func (g *WaitGroup) Error() error {
	return g.err
}
func (g *WaitGroup) Supervise(w wire.Wire, args ...interface{}) {
	go func(g *WaitGroup) {
		<-w.Context().Done()
		g.CancelWithError(w.Context().Err())
		g.sendCancelledMessage(w, args...)
	}(g)
	g.Wait()
	g.sendSuccessMessage(w, args...)
}
func (g *WaitGroup) sendCancelledMessage(w wire.Wire, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	sender := args[0].(msg.Sender)
	w.SendCancelled(sender, g.Error())
}
func (g *WaitGroup) sendSuccessMessage(w wire.Wire, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	sender := args[0].(msg.Sender)
	w.SendSuccess(sender)
}
