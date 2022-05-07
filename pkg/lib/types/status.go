package types

import "fmt"

type Status int

const (
	Running Status = iota
	Success
	Error
	Warning
	Cancelled
	Done
)

func (s Status) String() string {
	switch s {
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Cancelled:
		return "Cancelled"
	case Done:
		return "Done"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}
