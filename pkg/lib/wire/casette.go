package wire

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Casette []*msg.Feedback

func NewCasette() *Casette {
	return &Casette{}
}

type Memo struct {
	Errors        []*msg.Feedback
	Successes     []*msg.Feedback
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

func (c *Casette) process(m *msg.Feedback) {
	*c = append(*c, m)
}
func (c *Casette) Record(w *Wire) {
	for m := range w.GetFeedback() {
		c.process(m)
	}
}
func (c *Casette) RecordWithCallback(w *Wire, callback func(*msg.Feedback) error) error {
	for m := range w.GetFeedback() {
		c.process(m)
		if err := callback(m); err != nil {
			return err
		}
	}
	return nil
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
func (c *Casette) GetFeedbacks() []*msg.Feedback {
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
		Errors:    []*msg.Feedback{},
		Successes: []*msg.Feedback{},
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
				memo.ReadCount += m.Progress()
			} else if m.Sender == msg.OutputConnector {
				memo.WriteCount += m.Progress()
			}
		}
	}
	return memo
}
