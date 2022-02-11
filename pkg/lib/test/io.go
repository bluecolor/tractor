package test

import (
	"errors"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

func Record(w *wire.Wire) *wire.Casette {
	var inputSuccess, outputSuccess bool
	c := wire.NewCasette()

	cb := func(m *msg.Feedback) error {
		if m.Type == msg.Cancelled {
			return errors.New("cancelled")
		}
		if m.Type == msg.Error {
			return m.Content.(error)
		} else if m.Type == msg.Success {
			if m.IsInputSuccess() {
				inputSuccess = true
				println("closing data ....")
				w.CloseData()
			} else if m.IsOutputSuccess() {
				outputSuccess = true
			}
			if inputSuccess && outputSuccess {
				w.CloseFeedback()
			}
		}
		return nil
	}
	c.RecordWithCallback(w, cb)
	return c
}
