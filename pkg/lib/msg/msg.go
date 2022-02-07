package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type (
	MessageType int
	Sender      int
	Record      map[string]interface{}
)

const (
	Anonymous Sender = iota
	InputConnector
	OutputConnector
)
const (
	Data MessageType = iota
	Progress
	Success
	Error
	Info
	Warning
	Debug
	Cancelled
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

func (m MessageType) String() string {
	switch m {
	case Data:
		return "data"
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
	default:
		return fmt.Sprintf("%d", int(m))
	}
}

type Message struct {
	Type    MessageType
	Sender  Sender
	Content interface{}
}

func (m *Message) String() string {
	return fmt.Sprintf("%v: %v", m.Sender, m.Type)
}
func (m *Message) Data() []Record {
	return m.Content.([]Record)
}
func (m *Message) Count() int {
	if m.Type == Data {
		return len(m.Data())
	}
	return 0
}
func (m *Message) Error() error {
	return m.Content.(error)
}

func NewData(data interface{}, args ...interface{}) *Message {
	var content []Record
	var sender Sender = InputConnector
	switch val := data.(type) {
	case []Record:
		content = val
	case Record:
		content = []Record{val}
	default:
		return nil
	}
	if len(args) > 0 {
		sender = args[0].(Sender)
	}
	return &Message{
		Sender:  sender,
		Type:    Data,
		Content: content,
	}
}
func NewError(sender Sender, err error) *Message {
	return &Message{
		Type:    Error,
		Sender:  sender,
		Content: err,
	}
}
func NewSuccess(sender Sender, args ...interface{}) *Message {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Message{
		Type:    Success,
		Sender:  sender,
		Content: content,
	}
}
func NewInfo(sender Sender, content interface{}) *Message {
	return &Message{
		Type:    Info,
		Sender:  sender,
		Content: content,
	}
}
func NewWarning(sender Sender, content interface{}) *Message {
	return &Message{
		Type:    Warning,
		Sender:  sender,
		Content: content,
	}
}
func NewDebug(sender Sender, content interface{}) *Message {
	return &Message{
		Type:    Debug,
		Sender:  sender,
		Content: content,
	}
}
func NewOutputProgress(count int) *Message {
	return &Message{
		Sender:  OutputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewInputProgress(count int) *Message {
	return &Message{
		Sender:  InputConnector,
		Type:    Progress,
		Content: count,
	}
}
func NewProgress(sender Sender, count int) *Message {
	return &Message{
		Sender:  sender,
		Type:    Progress,
		Content: count,
	}
}
func NewCancelled(sender Sender, args ...interface{}) *Message {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Message{
		Type:    Cancelled,
		Sender:  sender,
		Content: content,
	}
}
func NewInputCancelled(args ...interface{}) *Message {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Message{
		Type:    Cancelled,
		Sender:  InputConnector,
		Content: content,
	}
}
func NewOutputCancelled(args ...interface{}) *Message {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Message{
		Type:    Cancelled,
		Sender:  OutputConnector,
		Content: content,
	}
}
