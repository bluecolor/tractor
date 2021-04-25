package tractor

import "github.com/bluecolor/tractor/config"

type Wire interface {
	SendMessage(message *Message)
	SendCatalog(catalog *config.Catalog)
	SendFeed(sender SenderType, feed *Feed)
	ReadMessages() <-chan *Message
	Close()
}
