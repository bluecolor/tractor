package esync

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
	"github.com/rs/zerolog/log"
)

type ManagedWaitGroup struct {
	*WaitGroup
	w   *wire.Wire
	ct  types.ConnectorType
	err error
}

func NewManagedWaitGroup(w *wire.Wire, ct types.ConnectorType) *ManagedWaitGroup {
	return &ManagedWaitGroup{
		WaitGroup: NewWaitGroup(),
		w:         w,
		ct:        ct,
	}
}

func (m *ManagedWaitGroup) Error() error {
	return m.err
}
func (m *ManagedWaitGroup) SetError(err error) {
	m.err = err
}

func (m *ManagedWaitGroup) CancelWithError(err error) {
	m.SetError(err)
	log.Debug().Msgf("Cancelling wait group for: %v", m.ct)
	m.WaitGroup.Cancel(m.ct)
}

func (m *ManagedWaitGroup) Wait() {
	go func(g *WaitGroup) {
		<-m.w.Context().Done()
		m.CancelWithError(m.w.Context().Err())
	}(m.WaitGroup)
	m.WaitGroup.Wait()
	log.Debug().Msgf("-------Wait group for: %v finished", m.ct)
	m.Finish()
}

func (m *ManagedWaitGroup) Finish() {
	if m.Error() == nil {
		m.w.SendSuccess(msg.SenderFromConnectorType(m.ct))
	}
}
