package test

import (
	"context"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func Record(w wire.Wire, cancel context.CancelFunc) *wire.Casette {
	var inputSuccess, outputSuccess bool

	cb := func(m *msg.Message, cancel context.CancelFunc) {
		if m.Type == msg.Error {
			cancel()
			w.Close()
		} else if m.Type == msg.Success {
			if m.Sender == msg.InputConnector {
				inputSuccess = true
				w.CloseData()
			} else if m.Sender == msg.OutputConnector {
				outputSuccess = true
			}
			if inputSuccess && outputSuccess {
				w.CloseFeedback()
			}
		}
	}
	c := wire.NewCasette()
	c.RecordWithCancellable(w, cancel, cb)
	return c
}
