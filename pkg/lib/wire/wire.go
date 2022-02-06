package wire

import (
	"context"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

type Status struct {
	inputError  error
	outputError error
	closed      bool
}

func (s *Status) HasError() bool {
	return s.inputError != nil || s.outputError != nil
}

type Wire struct {
	ctx      context.Context
	data     chan *msg.Message
	feedback chan *msg.Message
	status   *Status
}

func New(ctx context.Context) *Wire {
	return &Wire{
		ctx:      ctx,
		data:     make(chan *msg.Message, 1000),
		feedback: make(chan *msg.Message, 1000),
		status: &Status{
			closed: false,
		},
	}
}
func NewWithTimeout(timeout time.Duration) (*Wire, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return New(ctx), ctx, cancel
}

func (w *Wire) GetInputError() error {
	return w.status.inputError
}
func (w *Wire) GetOutputError() error {
	return w.status.outputError
}
func (w *Wire) IsClosed() bool {
	return w.status.closed
}
func (w *Wire) HasError() bool {
	return w.status.HasError()
}
func (w *Wire) SetError(sender msg.Sender, err error) {
	switch sender {
	case msg.InputConnector:
		w.status.inputError = err
	case msg.OutputConnector:
		w.status.outputError = err
	}
}
func (w *Wire) SetInputError(err error) {
	w.SetError(msg.InputConnector, err)
}
func (w *Wire) SetOutputError(err error) {
	w.SetError(msg.OutputConnector, err)
}
func (w *Wire) SetClosed() {
	w.status.closed = true
}
func (w *Wire) Context() context.Context {
	return w.ctx
}
func (w *Wire) Close() {
	w.CloseData()
	w.CloseFeedback()
	w.SetClosed()
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
func (w *Wire) SendOutputProgress(count int) {
	w.SendProgress(msg.OutputConnector, count)
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
