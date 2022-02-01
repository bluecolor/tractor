package wire

import feeds "github.com/bluecolor/tractor/pkg/lib/feeds"

type Wire struct {
	feedChannel          chan feeds.Feed
	errorChannel         chan feeds.Feed
	progressChannel      chan feeds.Feed
	readProgressChannel  chan feeds.Feed
	writeProgressChannel chan feeds.Feed
	dataChannel          chan feeds.Data
	readDoneChannel      chan bool
	writeDoneChannel     chan bool
	doneChannel          chan bool
}

func NewWire() Wire {
	return Wire{
		feedChannel:          make(chan feeds.Feed, 10000),
		errorChannel:         make(chan feeds.Feed, 10000),
		progressChannel:      make(chan feeds.Feed, 10000),
		readProgressChannel:  make(chan feeds.Feed, 10000),
		writeProgressChannel: make(chan feeds.Feed, 10000),
		dataChannel:          make(chan feeds.Data, 10000),
		readDoneChannel:      make(chan bool, 1),
		writeDoneChannel:     make(chan bool, 1),
	}
}
func (w *Wire) SignalDataTerm() {
	w.SendData(nil)
}
func (w *Wire) WriteWorkerDone() {
	w.SignalDataTerm()
}
func (w *Wire) ReadDone() {
	w.readDoneChannel <- true
	w.SignalDataTerm()
}
func (w *Wire) WriteDone() {
	w.writeDoneChannel <- true
}
func (w *Wire) Done() {
	w.doneChannel <- true
}
func (w *Wire) IsReadDone() chan bool {
	return w.readDoneChannel
}
func (w *Wire) IsWriteDone() chan bool {
	return w.writeDoneChannel
}
func (w *Wire) SendFeed(feed feeds.Feed) {
	switch feed.Type {
	case feeds.ErrorFeed:
		w.errorChannel <- feed
	case feeds.ProgressFeed:
		w.progressChannel <- feed
		switch feed.Sender {
		case feeds.SenderInputConnector:
			w.readProgressChannel <- feed
		case feeds.SenderOutputConnector:
			w.writeProgressChannel <- feed
		}
	}
	w.feedChannel <- feed
}
func (w *Wire) SendReadProgress(count int, args ...interface{}) {
	w.SendFeed(feeds.NewReadProgress(count, args...))
}
func (w *Wire) SendWriteProgress(count int, args ...interface{}) {
	w.SendFeed(feeds.NewWriteProgress(count, args...))
}
func (w *Wire) SendInputSuccessFeed() {
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderInputConnector))
}
func (w *Wire) SendOutputSuccessFeed() {
	w.SendFeed(feeds.NewSuccessFeed(feeds.SenderOutputConnector))
}
func (w *Wire) SendData(data feeds.Data) {
	w.dataChannel <- data
}
func (w *Wire) ReadData() <-chan feeds.Data {
	return w.dataChannel
}
func (w *Wire) ReadFeeds() <-chan feeds.Feed {
	return w.feedChannel
}
func (w *Wire) ReadErrorFeeds() <-chan feeds.Feed {
	return w.errorChannel
}
func (w *Wire) ReadProgressFeeds() <-chan feeds.Feed {
	return w.progressChannel
}
func (w *Wire) WriteProgressFeeds() <-chan feeds.Feed {
	return w.writeProgressChannel
}
func (w *Wire) CloseErrorFeed() {
	close(w.errorChannel)
}
func (w *Wire) CloseProgressFeed() {
	close(w.progressChannel)
}
func (w *Wire) CloseReadProgressFeed() {
	close(w.readProgressChannel)
}
func (w *Wire) CloseWriteProgressFeed() {
	close(w.writeProgressChannel)
}
func (w *Wire) CloseReadDone() {
	close(w.readDoneChannel)
}
func (w *Wire) CloseWriteDone() {
	close(w.writeDoneChannel)
}
func (w *Wire) Close() {
	// todo others
	w.CloseFeed()
	w.CloseData()
}
func (w *Wire) CloseData() {
	close(w.dataChannel)
}
func (w *Wire) CloseFeed() {
	close(w.feedChannel)
}
