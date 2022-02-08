package esync

// import (
// 	"context"
// 	"errors"

// 	"github.com/bluecolor/tractor/pkg/lib/msg"
// 	"github.com/bluecolor/tractor/pkg/lib/types"
// 	"github.com/bluecolor/tractor/pkg/lib/wire"
// )

// type ManagedWaitGroup struct {
// 	*WaitGroup
// 	w      *wire.Wire
// 	ct     types.ConnectorType
// 	err    error
// 	ctx    context.Context
// 	cancel context.CancelFunc
// }

// func NewManagedWaitGroup(w *wire.Wire, ct types.ConnectorType) *ManagedWaitGroup {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &ManagedWaitGroup{
// 		WaitGroup: NewWaitGroup(),
// 		w:         w,
// 		ct:        ct,
// 		ctx:       ctx,
// 		cancel:    cancel,
// 	}
// }

// func (m *ManagedWaitGroup) Context() context.Context {
// 	return m.ctx
// }
// func (m *ManagedWaitGroup) Error() error {
// 	return m.err
// }
// func (m *ManagedWaitGroup) SetError(err error) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	if m.err == nil || m.err != m.ctx.Err() {
// 		m.err = err
// 	}
// }
// func (m *ManagedWaitGroup) Cancel() {
// 	m.WaitGroup.Cancel()
// }
// func (m *ManagedWaitGroup) CancelWithError(err error) {
// 	m.SetError(err)
// 	m.Cancel()
// }
// func (m *ManagedWaitGroup) Wait() error {
// 	m.WaitGroup.Wait()
// 	return m.Finish()
// }
// func (m *ManagedWaitGroup) Finish() error {
// 	switch {
// 	case m.Error() == nil:
// 		m.w.SendSuccess(msg.SenderFromConnectorType(m.ct))
// 	case errors.Is(m.Error(), m.w.Context().Err()):
// 		m.w.SendCancelled(msg.SenderFromConnectorType(m.ct))
// 	default:
// 		m.w.SendError(msg.SenderFromConnectorType(m.ct), m.Error())
// 	}
// 	return m.Error()
// }
// func (m *ManagedWaitGroup) HandleError(err error) {
// 	m.SetError(err)
// 	m.cancel()
// }
