package esync

import (
	"context"
	"errors"
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type WaitGroup struct {
	mu        sync.Mutex
	wg        *sync.WaitGroup
	w         *wire.Wire
	ct        types.ConnectorType
	err       error
	ctx       context.Context
	cancel    context.CancelFunc
	cancelled bool
}

func NewWaitGroup(w *wire.Wire, ct types.ConnectorType) *WaitGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &WaitGroup{
		mu:     sync.Mutex{},
		wg:     &sync.WaitGroup{},
		w:      w,
		ct:     ct,
		ctx:    ctx,
		cancel: cancel,
	}
}
func (g *WaitGroup) Context() context.Context {
	return g.ctx
}
func (g *WaitGroup) Error() error {
	return g.err
}
func (g *WaitGroup) Cancelled() bool {
	return g.cancelled
}
func (g *WaitGroup) HandleError(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.err == nil {
		g.err = err
	}
	if err.Error() == "send on closed channel" {
		g.cancelled = true
	}
}
func (g *WaitGroup) Add(n int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.wg.Add(n)
}
func (g *WaitGroup) Done() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.wg.Done()
}
func (g *WaitGroup) Finish() error {
	sender := msg.SenderFromConnectorType(g.ct)
	defer g.w.SendDone(sender)
	if g.cancelled {
		g.w.SendCancelled(sender)
	} else if g.err != nil {
		g.w.SendError(sender, g.err)
	} else {
		g.w.SendSuccess(sender)
	}
	if g.cancelled {
		return nil
	}
	return g.err
}
func (wg *WaitGroup) Wait() error {
	wg.wg.Wait()
	return wg.Finish()
}
func (wg *WaitGroup) Cancel() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	wg.cancel()
	wg.cancelled = errors.Is(wg.err, wg.ctx.Err())
}
