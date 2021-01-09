package message

import (
	"encoding/json"

	"github.com/bluecolor/tractor/api/message/mt"
	"github.com/bluecolor/tractor/api/message/sender"
	"github.com/bluecolor/tractor/api/message/status"
	"github.com/bluecolor/tractor/api/schema"
)

// Message ...
type Message interface {
	MessageType() mt.Type
	Content() interface{}
	Sender() sender.Type
}

type message struct {
	sender      sender.Type
	messageType mt.Type
	content     []byte
}

type statusMessage struct {
	processedMessageCount uint64
	status                status.Type
	error                 error
}

// StatusMessage ...
type StatusMessage interface {
	Status() status.Type
}

func (m *message) Sender() sender.Type {
	return m.sender
}

func (m *message) Content() interface{} {
	switch m.messageType {
	case mt.Status:
		return json.Unmarshal(m.content, &statusMessage{})
	case mt.Schema:
		return json.Unmarshal(m.content, &schema.DataStore{})
	}
	return nil
}

func (m *statusMessage) Status() status.Type {
	return m.status
}
