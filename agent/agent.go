package agent

import (
	"github.com/bluecolor/tractor"
	_ "github.com/bluecolor/tractor/plugins/inputs/all"
)

type wire struct {
	Data chan *tractor.Data
	Feed chan *tractor.Feed
}

func NewWire(
	data chan *tractor.Data,
	feed chan *tractor.Feed,
) tractor.Wire {
	w := wire{
		Data: data,
		Feed: feed,
	}
	return &w
}

func (w *wire) SendData(data *tractor.Data) {
	w.Data <- data
}

func (w *wire) SendFeed(feed *tractor.Feed) {
	w.Feed <- feed
}
