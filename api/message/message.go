package message

// Type ...
type Type int

// StatusType ...
type StatusType Type

// SenderType ...
type SenderType Type

const (
	// Metadata ...
	Metadata Type = iota
	// Data ...
	Data
	// Status ...
	Status
)

const (
	// Running ...
	Running StatusType = iota
	// Error ...
	Error
	// Success ...
	Success
	// Done ...
	Done
)

const (
	// InputPlugin ...
	InputPlugin SenderType = iota
	// OutputPlugin ...
	OutputPlugin
)

// Message ...
type Message struct {
	Sender  SenderType
	Type    Type
	Content interface{}
}

// StatusInfo ...
type StatusInfo struct {
	ProcessedMessageCount uint64
	Status                StatusType
	Error                 error
}

// NewDataMessage ...
func NewDataMessage(content interface{}) Message {
	return Message{
		Sender:  InputPlugin,
		Type:    Data,
		Content: content,
	}
}
