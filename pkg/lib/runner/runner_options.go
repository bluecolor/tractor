package runner

import "github.com/bluecolor/tractor/pkg/lib/msg"

type OptionType int

const (
	SessionIDOption OptionType = iota
	FeedbackBackendOption
)

type Option struct {
	Type  OptionType
	Value interface{}
}

func WithSessionIDOption(sessionID string) Option {
	return Option{
		Type:  SessionIDOption,
		Value: sessionID,
	}
}
func WithFeedbackBackendOption(backend msg.FeedbackBackend) Option {
	return Option{
		Type:  FeedbackBackendOption,
		Value: backend,
	}
}
