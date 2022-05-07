package types

type ConnectorType int

const (
	Input ConnectorType = iota
	Output
)

func (c ConnectorType) String() string {
	switch c {
	case Input:
		return "Input"
	case Output:
		return "Output"
	}
	return ""
}
