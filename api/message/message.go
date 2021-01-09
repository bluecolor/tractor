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
	Content     []byte
}

// Status ...
type Status struct {
	ProcessedMessageCount uint64
	Status                status.Type
	Error                 error
}

// Data ...
type Data struct {
	records []byte
}
