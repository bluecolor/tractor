package runner

import "github.com/bluecolor/tractor/pkg/lib/msg"

type OptionType int

const (
	FeedBackendOption OptionType = iota
)

type Option struct {
	Type  OptionType
	Value interface{}
}

func WithFeedbackBackendOption(backend msg.FeedBackend) Option {
	return Option{
		Type:  FeedBackendOption,
		Value: backend,
	}
}

func GetFeedBackends(options ...Option) []msg.FeedBackend {
	var backends []msg.FeedBackend
	for _, opt := range options {
		if opt.Type == FeedBackendOption {
			backends = append(backends, opt.Value.(msg.FeedBackend))
		}
	}
	return backends
}
