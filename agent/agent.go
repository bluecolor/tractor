package agent

import (
	"github.com/bluecolor/tractor"
	_ "github.com/bluecolor/tractor/plugins/inputs/all"
	_ "github.com/bluecolor/tractor/plugins/outputs/all"
)

type wire struct {
	feedChannel chan tractor.Feed
	dataChannel chan tractor.Data
}

func NewWire() tractor.Wire {
	w := wire{
		feedChannel: make(chan tractor.Feed, 100),
		dataChannel: make(chan tractor.Data, 100),
	}
	return &w
}

func (w *wire) SendFeed(feed tractor.Feed) {
	w.feedChannel <- feed
}

func (w *wire) SendData(data tractor.Data) {
	w.dataChannel <- data
}

func (w *wire) ReadData() <-chan tractor.Data {
	return w.dataChannel
}

func (w *wire) ReadFeeds() <-chan tractor.Feed {
	return w.feedChannel
}

func (w *wire) Close() {
	close(w.feedChannel)
	close(w.dataChannel)
}

func (w *wire) CloseData() {
	close(w.dataChannel)
}

func (w *wire) CloseFeed() {
	close(w.feedChannel)
}
