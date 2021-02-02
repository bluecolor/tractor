package message

// Type ...
type Type int

// OrderType ...
type OrderType int

// SenderType ...
type SenderType int

const (
	// Error ...s
	Error Type = iota
	// Success ...
	Success
	// Progress ...
	Progress
	// Order ...
	Order
)

const (
	// Stop ...
	Stop OrderType = iota
)

const (
	// InputSupervisor ...
	InputSupervisor SenderType = iota
	// InputWorker ...
	InputWorker
	// InputPlugin ...
	InputPlugin
	// OutputSupervisor ...
	OutputSupervisor
	// OutputWorker ...
	OutputWorker
	// OutputPlugin ...
	OutputPlugin
	// Anonymous ...
	Anonymous
)

// Message ...
type Message struct {
	Sender  SenderType
	Type    Type
	Content interface{}
}

// NewErrorMessage ...
func NewErrorMessage(sender SenderType, args ...interface{}) *Message {
	return NewMessage(sender, Error, args)
}

// NewSuccessMessage ...
func NewSuccessMessage(sender SenderType, args ...interface{}) *Message {
	return NewMessage(sender, Success, args)
}

// NewStopOrder ...
func NewStopOrder(args ...interface{}) *Message {
	return NewOrder(Stop, args)
}

// NewOrder ...
func NewOrder(order OrderType, args ...interface{}) *Message {
	var sender SenderType = -1
	if len(args) > 0 {
		sender = args[0].(SenderType)
	}
	return &Message{
		Sender: sender, Type: Order, Content: order,
	}
}

// NewMessage ...
func NewMessage(sender SenderType, messageType Type, args ...interface{}) *Message {
	var content interface{}
	if len(args) > 0 {
		content = args[0]
	}
	return &Message{
		Sender: sender, Type: messageType, Content: content,
	}
}
