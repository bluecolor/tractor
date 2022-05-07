package wire

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/pairs"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/rs/zerolog/log"
)

type Wire interface {
	Close()
	CloseData()
	CloseFeeds()
	Send(message interface{}, args ...types.Pair)
	SendData(data interface{}, args ...types.Pair)
	ReceiveData() <-chan msg.Data
	SendFeed(feed msg.Feed)
	ReceiveFeed() <-chan msg.Feed
	SendSuccess(sender types.Sender)
	SendError(sender types.Sender, err error)
	SendProgress(types.Sender, int)
	SendInputProgress(progress int)
	SendOutputProgress(progress int)
	SendInputSuccess()
	SendInputError(err error)
	SendOutputSuccess()
	SendOutputError(err error)
}

type wire struct {
	data  chan msg.Data
	feeds chan msg.Feed
}

func New() Wire {
	return &wire{
		data:  make(chan msg.Data, 1000),
		feeds: make(chan msg.Feed, 100),
	}
}
func (w *wire) Close() {
	w.CloseData()
	w.CloseFeeds()
}
func (w *wire) CloseFeeds() {
	close(w.feeds)
}
func (w *wire) CloseData() {
	log.Debug().Msg("closing data channel....")
	close(w.data)
	log.Debug().Msg("data channel closed")
}
func (w *wire) Send(m interface{}, args ...types.Pair) {
	switch m.(type) {
	case msg.Data, msg.Record:
		w.SendData(m, args...)
	case msg.Feed:
		w.SendFeed(m.(msg.Feed))
	}
}
func (w *wire) SendData(data interface{}, args ...types.Pair) {
	message := msg.NewData(data)
	w.data <- message
	if pairs.GetOr(pairs.SendProgress, true, args...) {
		w.SendProgress(types.InputConnector, message.Count())
	}
}
func (w *wire) ReceiveData() <-chan msg.Data {
	return w.data
}
func (w *wire) SendFeed(feed msg.Feed) {
	w.feeds <- feed
}
func (w *wire) ReceiveFeed() <-chan msg.Feed {
	return w.feeds
}
func (w *wire) SendSuccess(sender types.Sender) {
	w.feeds <- msg.NewStatusFeed(sender, types.Success)
}
func (w *wire) SendError(sender types.Sender, err error) {
	w.feeds <- msg.NewStatusFeed(sender, types.Success, err)
}
func (w *wire) SendProgress(sender types.Sender, count int) {
	w.feeds <- msg.NewProgressFeed(sender, count)
}
func (w *wire) SendInputProgress(progress int) {
	w.SendProgress(types.InputConnector, progress)
}
func (w *wire) SendOutputProgress(progress int) {
	w.SendProgress(types.OutputConnector, progress)
}
func (w *wire) SendInputSuccess() {
	w.SendSuccess(types.InputConnector)
}
func (w *wire) SendInputError(err error) {
	w.SendError(types.InputConnector, err)
}
func (w *wire) SendOutputSuccess() {
	w.SendSuccess(types.OutputConnector)
}
func (w *wire) SendOutputError(err error) {
	w.SendError(types.OutputConnector, err)
}
func (w *wire) SendDone(s types.Sender) {
	w.feeds <- msg.NewStatusFeed(s, types.Done)
}
func (w *wire) SendInputDone() {
	w.SendDone(types.InputConnector)
}
func (w *wire) SendOutputDone() {
	w.SendDone(types.OutputConnector)
}
