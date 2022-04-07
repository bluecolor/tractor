package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type (
	SessionStatusType int
	FeedType          int
	Sender            int
	Record            []interface{}
	Data              []Record
)

const (
	Anonymous Sender = iota
	InputConnector
	OutputConnector
	Supervisor
	Driver
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

func (s Sender) String() string {
	switch s {
	case Anonymous:
		return "Anonymous"
	case InputConnector:
		return "InputConnector"
	case OutputConnector:
		return "OutputConnector"
	case Supervisor:
		return "Supervisor"
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
