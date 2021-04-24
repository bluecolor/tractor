package agent

import (
	"github.com/bluecolor/tractor"
	_ "github.com/bluecolor/tractor/plugins/inputs/all"
)

type wire struct {
	feedChannel    chan *tractor.Message
	messageChannel chan *tractor.Message
}

func NewWire() tractor.Wire {
	w := wire{
		feedChannel:    make(chan *tractor.Message, 100),
		messageChannel: make(chan *tractor.Message, 100),
	}
	return &w
}

func (w *wire) SendMessage(message *tractor.Message) {
	if message.Type == tractor.FeedMessage {
		w.feedChannel <- message
	} else {
		w.messageChannel <- message
	}
}

func (w *wire) ReadMessages() <-chan *tractor.Message {
	return w.messageChannel
}
