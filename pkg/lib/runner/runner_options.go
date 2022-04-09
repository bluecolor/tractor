package runner

import (
	"github.com/bluecolor/tractor/pkg/tasks/feedproc"
)

type OptionType int

const (
	FeedProcessorOption OptionType = iota
)

type Option struct {
	Type  OptionType
	Value interface{}
}

func WithFeedProcessorOption(processor *feedproc.FeedProcessor) Option {
	return Option{
		Type:  FeedProcessorOption,
		Value: processor,
	}
}
func GetFeedProcessor(options ...Option) *feedproc.FeedProcessor {
	for _, opt := range options {
		if opt.Type == FeedProcessorOption {
			return opt.Value.(*feedproc.FeedProcessor)
		}
	}
	return nil
}
