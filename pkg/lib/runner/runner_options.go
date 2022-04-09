package runner

import (
	"net/rpc"
)

type OptionType int

const (
	FeedClientOption OptionType = iota
)

type Option struct {
	Type  OptionType
	Value interface{}
}

func WithFeedClientOption(client *rpc.Client) Option {
	return Option{
		Type:  FeedClientOption,
		Value: client,
	}
}

func GetFeedClient(options ...Option) *rpc.Client {
	for _, opt := range options {
		if opt.Type == FeedClientOption {
			return opt.Value.(*rpc.Client)
		}
	}
	return nil
}
