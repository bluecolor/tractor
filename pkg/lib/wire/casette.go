package wire

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Casette []*msg.Feed

func NewCasette() *Casette {
	return &Casette{}
}

type Memo struct {
	errors           []*msg.Feed
	successes        []*msg.Feed
	hasInputError    bool
	hasInputSuccess  bool
	hasOutputError   bool
	hasOutputSuccess bool
	readCount        int
	writeCount       int
	feedbackCount    int
}

func (m *Memo) HasInputError() bool {
	for _, m := range m.errors {
		if m.Sender == msg.InputConnector {
			return true
		}
	}
	return false
}
func (m *Memo) HasOutputError() bool {
	for _, m := range m.errors {
		if m.Sender == msg.OutputConnector {
			return true
		}
	}
	return false
}
func (m *Memo) HasInputSuccess() bool {
	for _, m := range m.successes {
		if m.Sender == msg.InputConnector {
			return true
		}
	}
	return false
}
func (m *Memo) HasOutputSuccess() bool {
	for _, m := range m.successes {
		if m.Sender == msg.OutputConnector {
			return true
		}
	}
	return false
}
func (m *Memo) IsEmpty() bool {
	return m.feedbackCount == 0
}
func (m *Memo) HasError() bool {
	return len(m.errors) > 0
}
func (m *Memo) ReadCount() int {
	return m.readCount
}
func (m *Memo) WriteCount() int {
	return m.writeCount
}
func (m *Memo) Errors() []*msg.Feed {
	return m.errors
}
func (m *Memo) Successes() []*msg.Feed {
	return m.successes
}

func (c *Casette) process(m *msg.Feed) {
	*c = append(*c, m)
}
func (c *Casette) Record(w *Wire) {
	for m := range w.ReceiveFeedback() {
		c.process(m)
	}
}
func (c *Casette) RecordWithCallback(w *Wire, callback func(*msg.Feed) error) error {
	for m := range w.ReceiveFeedback() {
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
func (c *Casette) GetFeedbacks() []*msg.Feed {
	return *c
}
func (c *Casette) HasSuccess() bool {
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
func (c *Casette) HasError() bool {
	return !c.HasSuccess()
}
func (c *Casette) HasInputSuccess() bool {
	for _, m := range *c {
		if m.Sender == msg.InputConnector && m.Type == msg.Success {
			return true
		}
	}
	return false
}
func (c *Casette) HasOutputSuccess() bool {
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
func (c *Casette) Memo() *Memo {
	var memo *Memo = &Memo{
		errors:    []*msg.Feed{},
		successes: []*msg.Feed{},
	}
	if len(*c) == 0 {
		return memo
	}
	for _, m := range *c {
		memo.feedbackCount++
		switch m.Type {
		case msg.Error, msg.Cancelled:
			memo.errors = append(memo.errors, m)
			switch m.Sender {
			case msg.InputConnector:
				memo.hasInputError = true
			case msg.OutputConnector:
				memo.hasOutputError = true
			}
		case msg.Success:
			memo.successes = append(memo.successes, m)
			switch m.Sender {
			case msg.InputConnector:
				memo.hasInputSuccess = true
			case msg.OutputConnector:
				memo.hasOutputSuccess = true
			}
		case msg.Progress:
			if m.Sender == msg.InputConnector {
				memo.readCount += m.Progress()
			} else if m.Sender == msg.OutputConnector {
				memo.writeCount += m.Progress()
			}
		}
	}
	return memo
}
