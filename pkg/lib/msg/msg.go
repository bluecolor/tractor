package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type (
	FeedbackType int
	Sender       int
	Record       map[string]interface{}
	Data         []Record
)

const (
	Anonymous Sender = iota
	InputConnector
	OutputConnector
)
const (
	Progress FeedbackType = iota
	Success
	Error
	Info
	Warning
	Debug
	Cancelled
	Done
)

func SenderFromConnectorType(ct types.ConnectorType) Sender {
	switch ct {
	case types.InputConnector:
		return InputConnector
	case types.OutputConnector:
		return OutputConnector
	default:
		return Anonymous
	}
}

func (ft FeedbackType) String() string {
	switch ft {
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
func (s Sender) String() string {
	switch s {
	case Anonymous:
		return "Anonymous"
	case InputConnector:
		return "InputConnector"
	case OutputConnector:
		return "OutputConnector"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

func NewData(data interface{}, args ...interface{}) Data {
	var content []Record
	switch val := data.(type) {
	case []Record:
		content = val
	case Record:
		content = []Record{val}
	case Data:
		content = val
	default:
		return nil
	}
	return content
}

func (d Data) Count() int {
	return len(d)
}

type Feedback struct {
	Type    FeedbackType
	Sender  Sender
	Content interface{}
}

func (f *Feedback) String() string {
	return fmt.Sprintf("%v: %v", f.Sender, f.Type)
}
func (f *Feedback) Data() []Record {
	return f.Content.([]Record)
}

func (f *Feedback) Progress() int {
	switch f.Type {
	case Progress:
		return f.Content.(int)
	default:
		return 0
	}
}
func (f *Feedback) InputProgress() int {
	if f.Sender == InputConnector {
		return f.Progress()
	}
	return 0
}
func (f *Feedback) OutputProgress() int {
	if f.Sender == OutputConnector {
		return f.Progress()
	}
	return 0
}
func (f *Feedback) IsInputSuccess() bool {
	return f.Sender == InputConnector && f.Type == Success
}
func (f *Feedback) IsOutputSuccess() bool {
	return f.Sender == OutputConnector && f.Type == Success
}
func (f *Feedback) IsInputError() bool {
	return f.Sender == InputConnector && f.Type == Error
}
func (f *Feedback) IsOutputError() bool {
	return f.Sender == OutputConnector && f.Type == Error
}
func (f *Feedback) IsError() bool {
	return f.Type == Error
}
func (f *Feedback) IsSuccess() bool {
	return f.Type == Success
}
func (f *Feedback) IsCancelled() bool {
	return f.Type == Cancelled
}
func (f *Feedback) IsInputDone() bool {
	return f.Sender == InputConnector && f.Type == Done
}
func (f *Feedback) IsOutputDone() bool {
	return f.Sender == OutputConnector && f.Type == Done
}
func (f *Feedback) Error() error {
	return f.Content.(error)
}
func NewDone(sender Sender, args ...interface{}) *Feedback {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feedback{
		Type:    Done,
		Sender:  sender,
		Content: content,
	}
}
func NewError(sender Sender, err error) *Feedback {
	return &Feedback{
		Type:    Error,
		Sender:  sender,
		Content: err,
	}
}
func NewSuccess(sender Sender, args ...interface{}) *Feedback {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feedback{
		Type:    Success,
		Sender:  sender,
		Content: content,
	}
}
func NewInfo(sender Sender, content interface{}) *Feedback {
	return &Feedback{
		Type:    Info,
		Sender:  sender,
		Content: content,
	}
}
func NewWarning(sender Sender, content interface{}) *Feedback {
	return &Feedback{
		Type:    Warning,
		Sender:  sender,
		Content: content,
	}
}
func NewDebug(sender Sender, content interface{}) *Feedback {
	return &Feedback{
		Type:    Debug,
		Sender:  sender,
		Content: content,
	}
}
func NewOutputProgress(count int) *Feedback {
	return &Feedback{
		Sender:  OutputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewInputProgress(count int) *Feedback {
	return &Feedback{
		Sender:  InputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewProgress(sender Sender, count int) *Feedback {
	return &Feedback{
		Sender:  sender,
		Type:    Progress,
		Content: count,
	}
}
func NewCancelled(sender Sender, args ...interface{}) *Feedback {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feedback{
		Type:    Cancelled,
		Sender:  sender,
		Content: content,
	}
}
func NewInputCancelled(args ...interface{}) *Feedback {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feedback{
		Type:    Cancelled,
		Sender:  InputConnector,
		Content: content,
	}
}
func NewOutputCancelled(args ...interface{}) *Feedback {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Feedback{
		Type:    Cancelled,
		Sender:  OutputConnector,
		Content: content,
	}
}
