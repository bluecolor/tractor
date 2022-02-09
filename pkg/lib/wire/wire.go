package wire

import (
	"context"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Wire struct {
	data     chan msg.Data
	feedback chan *msg.Feedback
}

func New(ctx context.Context) *Wire {
	return &Wire{
		data:     make(chan msg.Data, 1000),
		feedback: make(chan *msg.Feedback, 100),
	}
}
func NewWithTimeout(timeout time.Duration) (*Wire, context.CancelFunc, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return New(ctx), cancel, ctx
}
func NewWithDefaultTimeout() (*Wire, context.CancelFunc, context.Context) {
	timeout := time.Second * 50000 //todo default timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return New(ctx), cancel, ctx
}

func (w *Wire) Close() {
	w.CloseData()
	w.CloseFeedback()
}
func (w *Wire) CloseFeedback() {
	close(w.feedback)
}
func (w *Wire) CloseData() {
	close(w.data)
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
func (w *Wire) GetData() <-chan msg.Data {
	return w.data
}
func (w *Wire) SendFeedback(feedback *msg.Feedback) {
	w.feedback <- feedback
}
func (w *Wire) GetFeedback() <-chan *msg.Feedback {
	return w.feedback
}
func (w *Wire) SendSuccess(sender msg.Sender, args ...interface{}) {
	w.feedback <- msg.NewSuccess(sender, args)
}
func (w *Wire) SendProgress(sender msg.Sender, count int) {
	w.feedback <- msg.NewProgress(sender, count)
}
func (w *Wire) SendError(sender msg.Sender, err error) {
	w.feedback <- msg.NewError(sender, err)
}
func (w *Wire) SendInfo(sender msg.Sender, content interface{}) {
	w.feedback <- msg.NewInfo(sender, content)
}
func (w *Wire) SendWarning(sender msg.Sender, content interface{}) {
	w.feedback <- msg.NewWarning(sender, content)
}
func (w *Wire) SendDebug(sender msg.Sender, content interface{}) {
	w.feedback <- msg.NewDebug(sender, content)
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
	w.feedback <- msg.NewCancelled(sender, args)
}
func (w *Wire) SendInputCancelled(args ...interface{}) {
	w.SendCancelled(msg.InputConnector, args...)
}
func (w *Wire) SendOutputCancelled(args ...interface{}) {
	w.SendCancelled(msg.OutputConnector, args...)
}
