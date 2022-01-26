package wire

import feeds "github.com/bluecolor/tractor/pkg/lib/feeds"

type Wire struct {
	FeedChannel chan feeds.Feed
	DataChannel chan feeds.Data
}

func NewWire() Wire {
	return Wire{
		FeedChannel: make(chan feeds.Feed, 10000),
		DataChannel: make(chan feeds.Data, 10000),
	}
}
func (w *Wire) SendFeed(feed feeds.Feed) {
	w.FeedChannel <- feed
}
func (w *Wire) SendData(data feeds.Data) {
	w.DataChannel <- data
}
func (w *Wire) ReadData() <-chan feeds.Data {
	return w.DataChannel
}
func (w *Wire) ReadFeed() <-chan feeds.Feed {
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
