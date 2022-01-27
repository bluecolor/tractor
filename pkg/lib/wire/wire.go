package wire

import feeds "github.com/bluecolor/tractor/pkg/lib/feeds"

type Wire struct {
	FeedChannel      chan feeds.Feed
	DataChannel      chan feeds.Data
	ReadDoneChannel  chan bool
	WriteDoneChannel chan bool
	DoneChannel      chan bool
}

func NewWire() Wire {
	return Wire{
		FeedChannel:      make(chan feeds.Feed, 10000),
		DataChannel:      make(chan feeds.Data, 10000),
		ReadDoneChannel:  make(chan bool, 1),
		WriteDoneChannel: make(chan bool, 1),
	}
}
func (w *Wire) ReadDone() {
	w.ReadDoneChannel <- true
}
func (w *Wire) WriteDone() {
	w.WriteDoneChannel <- true
}
func (w *Wire) Done() {
	w.DoneChannel <- true
}
func (w *Wire) IsReadDone() chan bool {
	return w.ReadDoneChannel
}
func (w *Wire) IsWriteDone() chan bool {
	return w.WriteDoneChannel
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
