package message

import (
	"github.com/bluecolor/tractor/api/message/mt"
	"github.com/bluecolor/tractor/api/message/sender"
	"github.com/bluecolor/tractor/api/message/status"
)

// Message ...
type Message struct {
	Sender      sender.Type
	MessageType mt.Type
	Content     interface{}
}

// Status ...
type Status struct {
	ProcessedMessageCount uint64
	Status                status.Type
	Error                 error
}

// Data ...
type Data struct {
	Content [][]interface{}
}

// NewDataMessage ...
func NewDataMessage(content interface{}) Message {
	return Message{
		Sender:      sender.InputPlugin,
		MessageType: mt.Data,
		Content:     content,
	}
}
