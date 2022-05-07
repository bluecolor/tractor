package types

import "fmt"

type Sender int

const (
	Anonymous Sender = iota
	InputConnector
	OutputConnector
	Driver
)

func (s Sender) String() string {
	switch s {
	case Anonymous:
		return "Anonymous"
	case InputConnector:
		return "InputConnector"
	case OutputConnector:
		return "OutputConnector"
	case Driver:
		return "Driver"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}
