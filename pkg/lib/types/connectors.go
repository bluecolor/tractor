package types

type ConnectorType int

const (
	InputConnector ConnectorType = iota
	OutputConnector
)

func (c ConnectorType) String() string {
	switch c {
	case InputConnector:
		return "InputConnector"
	case OutputConnector:
		return "OutputConnector"
	}
	return ""
}
