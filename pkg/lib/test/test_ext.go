package test

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Expect struct {
	TestReadCount  bool
	ReadCount      int
	TestWriteCount bool
	WriteCount     int
}

type Result struct {
	ReadCount  int
	WriteCount int
	Feedbacks  []*msg.Message
}

func (r *Result) ProcessFeedback(f *msg.Message) {
	r.Feedbacks = append(r.Feedbacks, f)
	if f.Sender == msg.InputConnector {
		switch f.Type {
		case msg.Progress:
			r.ReadCount += f.Content.(int)


	return nil
}
func NewResult() *Result {
	return &Result{
		Feedbacks: make([]*msg.Message, 0),
	}
}

func TestExt(w wire.Wire, e Expect) error {

	r := NewResult()
	for f := range w.GetFeedback() {
		r.ProcessFeedback(f)
	}

	return nil
}
