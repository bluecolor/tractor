package wire

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/rs/zerolog/log"
)

type Wire struct {
	data  chan msg.Data
	feeds chan *msg.Feed
}

func New() *Wire {
	return &Wire{
		data:  make(chan msg.Data, 1000),
		feeds: make(chan *msg.Feed, 100),
	}
}

func (w *Wire) Close() {
	w.CloseData()
	w.CloseFeeds()
}
func (w *Wire) CloseFeeds() {
	close(w.feeds)
}
func (w *Wire) CloseData() {
	log.Debug().Msg("closing data channel....")
	close(w.data)
	log.Debug().Msg("data channel closed")
}
func (w *Wire) SendData(data interface{}, args ...interface{}) {
	var sendProgress bool = true
	if len(args) > 0 {
		sendProgress = args[0].(bool)
	}
	message := msg.NewData(data)
	w.data <- message
	if sendProgress {
		w.SendInputProgress(message.Count())
	}
}
func (w *Wire) ReceiveData() <-chan msg.Data {
	return w.data
}
func (w *Wire) SendFeed(feed *msg.Feed) {
	w.feeds <- feed
}
func (w *Wire) ReceiveFeedback() <-chan *msg.Feed {
	return w.feeds
}
func (w *Wire) SendSuccess(sender msg.Sender, args ...interface{}) {
	w.feeds <- msg.NewSuccess(sender, args)
}
func (w *Wire) SendProgress(sender msg.Sender, count int) {
	w.feeds <- msg.NewProgress(sender, count)
}
func (w *Wire) SendError(sender msg.Sender, err error) {
	w.feeds <- msg.NewError(sender, err)
}
func (w *Wire) SendInfo(sender msg.Sender, content interface{}) {
	w.feeds <- msg.NewInfo(sender, content)
}
func (w *Wire) SendWarning(sender msg.Sender, content interface{}) {
	w.feeds <- msg.NewWarning(sender, content)
}
func (w *Wire) SendDebug(sender msg.Sender, content interface{}) {
	w.feeds <- msg.NewDebug(sender, content)
}
func (w *Wire) SendInputProgress(progress int) {
	w.SendProgress(msg.InputConnector, progress)
}
func (w *Wire) SendInputSuccess(args ...interface{}) {
	w.SendSuccess(msg.InputConnector, args...)
}
func (w *Wire) SendInputError(err error) {
	w.SendError(msg.InputConnector, err)
}
func (w *Wire) SendInputInfo(content interface{}) {
	w.SendInfo(msg.InputConnector, content)
}
func (w *Wire) SendInputWarning(content interface{}) {
	w.SendWarning(msg.InputConnector, content)
}
func (w *Wire) SendInputDebug(content interface{}) {
	w.SendDebug(msg.InputConnector, content)
}
func (w *Wire) SendOutputProgress(progress int) {
	w.SendProgress(msg.OutputConnector, progress)
}
func (w *Wire) SendOutputSuccess(args ...interface{}) {
	w.SendSuccess(msg.OutputConnector, args...)
}
func (w *Wire) SendOutputError(err error) {
	w.SendError(msg.OutputConnector, err)
}
func (w *Wire) SendOutputInfo(content interface{}) {
	w.SendInfo(msg.OutputConnector, content)
}
func (w *Wire) SendOutputWarning(content interface{}) {
	w.SendWarning(msg.OutputConnector, content)
}
func (w *Wire) SendOutputDebug(content interface{}) {
	w.SendDebug(msg.OutputConnector, content)
}
func (w *Wire) SendCancelled(sender msg.Sender, args ...interface{}) {
	w.feeds <- msg.NewCancelled(sender, args)
}
func (w *Wire) SendInputCancelled(args ...interface{}) {
	w.SendCancelled(msg.InputConnector, args...)
}
func (w *Wire) SendOutputCancelled(args ...interface{}) {
	w.SendCancelled(msg.OutputConnector, args...)
}
func (w *Wire) SendDone(sender msg.Sender, args ...interface{}) {
	w.feeds <- msg.NewDone(sender, args)
}
func (w *Wire) SendInputDone(args ...interface{}) {
	w.SendDone(msg.InputConnector, args...)
}
func (w *Wire) SendOutputDone(args ...interface{}) {
	w.SendDone(msg.OutputConnector, args...)
}
