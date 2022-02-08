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
func (wg *WaitGroup) Context() context.Context {
	return wg.ctx
}
func (wg *WaitGroup) Error() error {
	return wg.err
}
func (wg *WaitGroup) Cancelled() bool {
	return wg.cancelled
}
func (wg *WaitGroup) HandleError(err error) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if errors.Is(err, wg.w.Context().Err()) {
		wg.cancelled = true
	} else if !errors.Is(err, wg.ctx.Err()) {
		wg.err = err
	}
}
func (wg *WaitGroup) Add(n int) {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	wg.wg.Add(n)
}
func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	defer wg.mu.Unlock()
	wg.wg.Done()
}
func (wg *WaitGroup) Finish() error {
	sender := msg.SenderFromConnectorType(wg.ct)
	if wg.cancelled {
		wg.w.SendCancelled(sender)
	} else if wg.err != nil {
		wg.w.SendError(sender, wg.err)
	} else {
		wg.w.SendSuccess(sender)
	}
	return wg.err
}
func (wg *WaitGroup) Wait() error {
	wg.wg.Wait()
	return wg.Finish()
}
