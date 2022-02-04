package wire

import feeds "github.com/bluecolor/tractor/pkg/lib/feeds"

type Wire struct {
	// data
	dataChannel chan feeds.Data
	// feeds
	feedChannel         chan feeds.Feed // all feeds
	errorFeedChannel    chan feeds.Feed // error feeds
	successFeedChannel  chan feeds.Feed // success feeds
	progressFeedChannel chan feeds.Feed // progress feeds
	// read
	readProgressFeedChannel chan feeds.Feed // read progress feeds
	readErrorFeedChannel    chan feeds.Feed // read error feeds
	readSuccessFeedChannel  chan feeds.Feed // read success feeds
	// write
	writeProgressFeedChannel chan feeds.Feed // write progress feeds
	writeErrorFeedChannel    chan feeds.Feed // write error feeds
	writeSuccessFeedChannel  chan feeds.Feed // write success feeds
	// read
	isReadSuccessChannel chan bool
	// write
	isWriteSuccessChannel chan bool // write
	isWriteErrorChannel   chan bool
	isWriteDoneChannel    chan bool
	// read/write done
	isDoneChannel chan bool
}

func New() Wire {
	return Wire{
		dataChannel:              make(chan feeds.Data, 1000),
		feedChannel:              make(chan feeds.Feed, 1000),
		errorFeedChannel:         make(chan feeds.Feed, 1000),
		progressFeedChannel:      make(chan feeds.Feed, 1000),
		readProgressFeedChannel:  make(chan feeds.Feed, 1000),
		readErrorFeedChannel:     make(chan feeds.Feed, 1000),
		readSuccessFeedChannel:   make(chan feeds.Feed, 1000),
		writeProgressFeedChannel: make(chan feeds.Feed, 1000),
		writeErrorFeedChannel:    make(chan feeds.Feed, 1000),
		writeSuccessFeedChannel:  make(chan feeds.Feed, 1000),
		isReadSuccessChannel:     make(chan bool, 1),
		isReadErrorChannel:       make(chan bool, 1),
		isReadDoneChannel:        make(chan bool, 1),
		isWriteSuccessChannel:    make(chan bool, 1),
		isWriteErrorChannel:      make(chan bool, 1),
		isWriteDoneChannel:       make(chan bool, 1),
		isDoneChannel:            make(chan bool, 1),
	}
}

//  parallel read/write workers
// -------------------------------------------------------------------
func (w *Wire) ReadWorkerDone() {
	// done message for parallel read workers
}
func (w *Wire) WriteWorkerDone() {
	// done message for parallel write workers
}

// read/write success
// -------------------------------------------------------------------
func (w *Wire) ReadSuccess() {
	// all read workers done with success
	w.isReadSuccessChannel <- true
}
func (w *Wire) WriteSuccess() {
	// all write workers done with success
	w.writeSuccessChannel <- true
}

// read/write done
// -------------------------------------------------------------------
func (w *Wire) ReadDone() {
	// all read workers done whether success or not
	w.readDoneChannel <- true
}
func (w *Wire) WriteDone() {
	w.writeDoneChannel <- true
}

func (w *Wire) Done() {
	w.doneChannel <- true
}
func (w *Wire) IsDone() chan bool {
	return w.doneChannel
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
func (w *Wire) SendErrorFeed(sender feeds.SenderType, err error) {
	w.SendFeed(feeds.NewError(sender, err))
}
func (w *Wire) SendReadErrorFeed(err error) {
	w.SendFeed(feeds.NewError(feeds.SenderInputConnector, err))
}
func (w *Wire) SendWriteErrorFeed(err error) {
	w.SendFeed(feeds.NewError(feeds.SenderOutputConnector, err))
}
func (w *Wire) SendInputErrorFeed(err error) {
	w.SendFeed(feeds.NewError(feeds.SenderInputConnector, err))
}
func (w *Wire) SendOutputErrorFeed(err error) {
	w.SendFeed(feeds.NewError(feeds.SenderOutputConnector, err))
}
func (w *Wire) SendReadProgress(count int, args ...interface{}) {
	w.SendFeed(feeds.NewReadProgress(count, args...))
}
func (w *Wire) SendWriteProgress(count int, args ...interface{}) {
	w.SendFeed(feeds.NewWriteProgress(count, args...))
}
func (w *Wire) SendInputSuccessFeed() {
	w.SendFeed(feeds.NewSuccess(feeds.SenderInputConnector))
}
func (w *Wire) SendOutputSuccessFeed() {
	w.SendFeed(feeds.NewSuccess(feeds.SenderOutputConnector))
}
func (w *Wire) SendReadSuccessFeed() {
	w.SendInputSuccessFeed()
}
func (w *Wire) SendWriteSuccessFeed() {
	w.SendOutputSuccessFeed()
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
func (w *Wire) CloseReadProgress() {
	close(w.readProgressChannel)
}
func (w *Wire) CloseWriteProgress() {
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
