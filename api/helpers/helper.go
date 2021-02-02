package helper

import (
	"errors"

	"github.com/bluecolor/tractor/api/helpers/message"
)

// Supervisor ...
type Supervisor struct{}

// Supervise ...
func (s *Supervisor) Supervise(in chan *message.Message, outs []chan *message.Message) error {
	var successCount int = 0
	for m := range in {
		if m.Type == message.Success {
			if successCount == len(outs) {
				break
			}
		} else if m.Type == message.Error {
			order := message.NewStopOrder()
			for _, o := range outs {
				o <- order
			}
			return errors.New("One of the child sessions failed")
		}
	}
	if successCount == len(outs) {
		return nil
	}
	return errors.New("Can not get success message from all childs")
}
