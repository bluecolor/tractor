package wire

import (
	"context"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Casette []*msg.Message

func NewCasette() *Casette {
	return &Casette{}
}

type Memo struct {
	Errors        []*msg.Message
	Successes     []*msg.Message
	ReadCount     int
	WriteCount    int
	FeedbackCount int
}

func (m *Memo) IsEmpty() bool {
	return m.FeedbackCount == 0
}
func (m *Memo) HasError() bool {
	return len(m.Errors) > 0
}

func (c *Casette) process(m *msg.Message) {
	*c = append(*c, m)
}
func (c *Casette) Record(w *Wire) {
	for m := range w.GetFeedback() {
		c.process(m)
	}
}
func (c *Casette) RecordWithCallback(w *Wire, callback func(*msg.Message)) {
	for m := range w.GetFeedback() {
		c.process(m)
		callback(m)
	}
}
func (c *Casette) RecordWithCancellable(w *Wire, cancel context.CancelFunc, callback func(*msg.Message) error) {
	for {
		select {
		case m, ok := <-w.GetFeedback():
			if !ok {
				return
			}
			c.process(m)
			if err := callback(m); err != nil {
				return
			}
		case <-w.Context().Done():
			return
		}
	}
}
func (c *Casette) GetReadCount() (count int) {
	for _, m := range *c {
		if m.Sender == msg.InputConnector && m.Type == msg.Progress {
			count += m.Content.(int)
		}
	}
	return
}
func (c *Casette) GetWriteCount() (count int) {
	for _, m := range *c {
		if m.Sender == msg.OutputConnector && m.Type == msg.Progress {
			count += m.Content.(int)
		}
	}
	return
}
func (c *Casette) GetFeedbacks() []*msg.Message {
	return *c
}
func (c *Casette) IsSuccess() bool {
	var readSuccess, writeSuccess bool
	for _, m := range *c {
		if m.Sender == msg.InputConnector && m.Type == msg.Success {
			readSuccess = true
		}
		if m.Sender == msg.OutputConnector && m.Type == msg.Success {
			writeSuccess = true
		}
		if readSuccess && writeSuccess {
			return true
		}
	}
	return false
}
func (c *Casette) IsError() bool {
	return !c.IsSuccess()
}
func (c *Casette) IsReadSuccess() bool {
	for _, m := range *c {
		if m.Sender == msg.InputConnector && m.Type == msg.Success {
			return true
		}
	}
	return false
}
func (c *Casette) IsWriteSuccess() bool {
	for _, m := range *c {
		if m.Sender == msg.OutputConnector && m.Type == msg.Success {
			return true
		}
	}
	return false
}
func (c *Casette) Errors() []error {
	var errors []error
	for _, m := range *c {
		if m.Type == msg.Error || m.Type == msg.Cancelled {
			errors = append(errors, m.Content.(error))
		}
	}
	return errors
}
func (c *Casette) GetMemo() *Memo {
	var memo *Memo = &Memo{
		Errors:    []*msg.Message{},
		Successes: []*msg.Message{},
	}
	if len(*c) == 0 {
		return memo
	}
	for _, m := range *c {
		memo.FeedbackCount++
		switch m.Type {
		case msg.Error, msg.Cancelled:
			memo.Errors = append(memo.Errors, m)
		case msg.Success:
			memo.Successes = append(memo.Successes, m)
		case msg.Progress:
			if m.Sender == msg.InputConnector {
				memo.ReadCount += m.Count()
			} else if m.Sender == msg.OutputConnector {
				memo.WriteCount += m.Count()
			}
		}
	}
	return memo
}
