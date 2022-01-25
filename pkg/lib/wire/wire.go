package wire

import "github.com/bluecolor/tractor/pkg/lib/feed"

type Wire struct {
	FeedChannel chan feed.Feed
	DataChannel chan feed.Data
}

func NewWire() *Wire {
	return &Wire{
		FeedChannel: make(chan feed.Feed, 10000),
		DataChannel: make(chan feed.Data, 10000),
	}
}
func (w *Wire) SendFeed(feed feed.Feed) {
	w.FeedChannel <- feed
}
func (w *Wire) SendData(data feed.Data) {
	w.DataChannel <- data
}
func (w *Wire) ReadData() <-chan feed.Data {
	return w.DataChannel
}
func (w *Wire) ReadFeed() <-chan feed.Feed {
	return w.FeedChannel
}
func (w *Wire) Close() {
	close(w.FeedChannel)
	close(w.DataChannel)
}
func (w *Wire) CloseData() {
	close(w.DataChannel)
}
func (w *Wire) CloseFeed() {
	close(w.FeedChannel)
}
