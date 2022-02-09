package test

import (
	"context"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func Record(w *wire.Wire, cancel context.CancelFunc) *wire.Casette {
	var inputSuccess, outputSuccess bool

	cb := func(m *msg.Feedback) error {
		if m.Type == msg.Error {
			cancel()
			return m.Content.(error)
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
		return nil
	}
	c := wire.NewCasette()
	c.RecordWithCallback(w, cb)
	return c
}
