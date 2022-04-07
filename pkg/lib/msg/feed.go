package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

const (
	Running FeedType = iota
	Progress
	Success
	Error
	Info
	Warning
	Debug
	Cancelled
	Done
)

func (ft FeedType) String() string {
	switch ft {
	case Running:
		return "Running"
	case Progress:
		return "Progress"
	case Success:
		return "Success"
	case Error:
		return "Error"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Debug:
		return "Debug"
	case Cancelled:
		return "Cancelled"
	case Done:
		return "Done"
	default:
		return fmt.Sprintf("%d", int(ft))
	}
}

type FeedBackend interface {
	Process(sessionID string, feed *Feed) error
}
type Feed struct {
	Type    FeedType    `json:"type"`
	Sender  Sender      `json:"sender"`
	Content interface{} `json:"content"`
}

func (f *Feed) String() string {
	return fmt.Sprintf("%v: %v", f.Sender, f.Type)
}
func (f *Feed) Data() []Record {
	return f.Content.([]Record)
}
func (f *Feed) Progress() int {
	switch f.Type {
	case Progress:
		return f.Content.(int)
	default:
		return 0
	}
}
func (f *Feed) InputProgress() int {
	if f.Sender == InputConnector {
		return f.Progress()
	}
	return 0
}
func (f *Feed) OutputProgress() int {
	if f.Sender == OutputConnector {
		return f.Progress()
	}
	return 0
}

func (f *Feed) IsProgress() bool {
	return f.Type == Progress
}
func (f *Feed) IsStatus() bool {
	return f.Type == Running ||
		f.Type == Done ||
		f.Type == Error ||
		f.Type == Success ||
		f.Type == Warning ||
		f.Type == Cancelled
}
func (f *Feed) IsSessionStatus() bool {
	return f.IsStatus() && f.Sender == Driver
}
func (f *Feed) IsIOStatus() bool {
	return (f.Sender == InputConnector || f.Sender == OutputConnector) && f.IsStatus()
}
func (f *Feed) IsSessionRunning() bool {
	return f.Type == Running && f.Sender == Driver
}
func (f *Feed) IsSessionSuccess() bool {
	return f.Type == Success && f.Sender == Driver
}
func (f *Feed) IsSessionError() bool {
	return f.Type == Error && f.Sender == Driver
}
func (f *Feed) IsInputSuccess() bool {
	return f.Sender == InputConnector && f.Type == Success
}
func (f *Feed) IsOutputSuccess() bool {
	return f.Sender == OutputConnector && f.Type == Success
}
func (f *Feed) IsInputError() bool {
	return f.Sender == InputConnector && f.Type == Error
}
func (f *Feed) IsOutputError() bool {
	return f.Sender == OutputConnector && f.Type == Error
}
func (f *Feed) IsError() bool {
	return f.Type == Error
}
func (f *Feed) IsSuccess() bool {
	return f.Type == Success
}
func (f *Feed) IsRunning() bool {
	return f.Type == Running
}
func (f *Feed) IsDone() bool {
	return f.Type == Done
}
func (f *Feed) IsCancelled() bool {
	return f.Type == Cancelled
}
func (f *Feed) IsInputDone() bool {
	return f.Sender == InputConnector && f.Type == Done
}
func (f *Feed) IsOutputDone() bool {
	return f.Sender == OutputConnector && f.Type == Done
}

func (f *Feed) Error() error {
	return f.Content.(error)
}
func (f *Feed) ErrorWithSource() (types.ErrorSource, error) {
	return f.ErrorSource(), f.Content.(error)
}
func (f *Feed) ErrorSource() types.ErrorSource {
	switch f.Sender {
	case InputConnector:
		return types.InputError
	case OutputConnector:
		return types.OutputError
	case Supervisor:
		return types.SupervisorError
	}
	return types.UnknownErrorSource
}

func NewDone(sender Sender, args ...interface{}) *Feed {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feed{
		Type:    Done,
		Sender:  sender,
		Content: content,
	}
}
func NewError(sender Sender, err error) *Feed {
	return &Feed{
		Type:    Error,
		Sender:  sender,
		Content: err,
	}
}
func NewSuccess(sender Sender, args ...interface{}) *Feed {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feed{
		Type:    Success,
		Sender:  sender,
		Content: content,
	}
}
func NewInfo(sender Sender, content interface{}) *Feed {
	return &Feed{
		Type:    Info,
		Sender:  sender,
		Content: content,
	}
}
func NewSessionRunning() *Feed {
	return &Feed{
		Sender: Driver,
		Type:   Running,
	}
}
func NewSessionError() *Feed {
	return &Feed{
		Sender: Driver,
		Type:   Error,
	}
}
func NewSessionDone() *Feed {
	return &Feed{
		Sender: Driver,
		Type:   Done,
	}
}
func NewSessionSuccess() *Feed {
	return &Feed{
		Sender: Driver,
		Type:   Success,
	}
}
func NewWarning(sender Sender, content interface{}) *Feed {
	return &Feed{
		Type:    Warning,
		Sender:  sender,
		Content: content,
	}
}
func NewDebug(sender Sender, content interface{}) *Feed {
	return &Feed{
		Type:    Debug,
		Sender:  sender,
		Content: content,
	}
}
func NewOutputProgress(count int) *Feed {
	return &Feed{
		Sender:  OutputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewInputProgress(count int) *Feed {
	return &Feed{
		Sender:  InputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewProgress(sender Sender, count int) *Feed {
	return &Feed{
		Sender:  sender,
		Type:    Progress,
		Content: count,
	}
}
func NewCancelled(sender Sender, args ...interface{}) *Feed {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feed{
		Type:    Cancelled,
		Sender:  sender,
		Content: content,
	}
}
func NewInputCancelled(args ...interface{}) *Feed {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feed{
		Type:    Cancelled,
		Sender:  InputConnector,
		Content: content,
	}
}
func NewOutputCancelled(args ...interface{}) *Feed {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feed{
		Type:    Cancelled,
		Sender:  OutputConnector,
		Content: content,
	}
}
