package wire

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Wire struct {
	data     chan *msg.Message
	feedback chan *msg.Message
}

func New() Wire {
	return Wire{
		data:     make(chan *msg.Message, 1000),
		feedback: make(chan *msg.Message, 1000),
	}
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
func (w *Wire) SendData(data interface{}) {
	w.data <- msg.NewData(data)
}
func (w *Wire) GetData() <-chan *msg.Message {
	return w.data
}
func (w *Wire) GetFeedback() <-chan *msg.Message {
	return w.feedback
}
func (w *Wire) SendSuccess(sender msg.Sender, args ...interface{}) {
	w.feedback <- msg.NewSuccess(sender, args)
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
	w.feedback <- msg.NewInputProgress(progress)
}
func (w *Wire) SendInputSuccess(args ...interface{}) {
	sender := msg.InputConnector
	w.feedback <- msg.NewSuccess(sender, args)
}
func (w *Wire) SendInputError(err error) {
	sender := msg.InputConnector
	w.feedback <- msg.NewError(sender, err)
}
func (w *Wire) SendInputInfo(content interface{}) {
	sender := msg.InputConnector
	w.feedback <- msg.NewInfo(sender, content)
}
func (w *Wire) SendInputWarning(content interface{}) {
	sender := msg.InputConnector
	w.feedback <- msg.NewWarning(sender, content)
}
func (w *Wire) SendInputDebug(content interface{}) {
	sender := msg.InputConnector
	w.feedback <- msg.NewDebug(sender, content)
}
func (w *Wire) SendOutputProgress(count int) {
	w.feedback <- msg.NewOutputProgress(count)
}
func (w *Wire) SendOutputSuccess(args ...interface{}) {
	sender := msg.OutputConnector
	w.feedback <- msg.NewSuccess(sender, args)
}
func (w *Wire) SendOutputError(err error) {
	sender := msg.OutputConnector
	w.feedback <- msg.NewError(sender, err)
}
func (w *Wire) SendOutputInfo(content interface{}) {
	sender := msg.OutputConnector
	w.feedback <- msg.NewInfo(sender, content)
}
func (w *Wire) SendOutputWarning(content interface{}) {
	sender := msg.OutputConnector
	w.feedback <- msg.NewWarning(sender, content)
}
func (w *Wire) SendOutputDebug(content interface{}) {
	sender := msg.OutputConnector
	w.feedback <- msg.NewDebug(sender, content)
}
