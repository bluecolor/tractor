package agent

import (
	"github.com/bluecolor/tractor"
	"github.com/bluecolor/tractor/config"
	_ "github.com/bluecolor/tractor/plugins/inputs/all"
	_ "github.com/bluecolor/tractor/plugins/outputs/all"
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

func (w *wire) SendCatalog(catalog *config.Catalog) {
	message := tractor.NewCatalogMessage(catalog)
	w.messageChannel <- message
}

func (w *wire) SendFeed(sender tractor.SenderType, feed *tractor.Feed) {
	message := tractor.NewFeed(sender, feed)
	w.feedChannel <- message
}

func (w *wire) ReadMessages() <-chan *tractor.Message {
	return w.messageChannel
}

func (w *wire) Close() {
	close(w.feedChannel)
	close(w.messageChannel)
}
