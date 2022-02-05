package wire

import "github.com/bluecolor/tractor/pkg/lib/msg"

type Casette []*msg.Message

func NewCasette() *Casette {
	return &Casette{}
}

func (c *Casette) process(m *msg.Message) {
	*c = append(*c, m)
}
func (c *Casette) Record(w *Wire) {
	for m := range w.GetFeedback() {
		c.process(m)
	}
}
func (c *Casette) RecordWithCallback(w Wire, callback func(*msg.Message)) {
	for m := range w.GetFeedback() {
		callback(m)
		c.process(m)
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
