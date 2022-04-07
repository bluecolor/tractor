package feedbackend

import (
	"net"
	"net/rpc"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog/log"
)

type FeedBackend struct {
	config conf.FeedBackend
}

func New(config conf.FeedBackend) *FeedBackend {
	return &FeedBackend{
		config: config,
	}
}

func (f *FeedBackend) Start(exit chan error) {
	addr, err := net.ResolveTCPAddr("tcp", f.config.Addr)
	if err != nil {
		log.Error().Err(err).Msg("failed to resolve tcp address")
		exit <- err
	}
	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen tcp")
		exit <- err
	}
	r, err := repo.New(f.config.DB)
	if err != nil {
		log.Error().Err(err).Msg("failed to create repo")
		exit <- err
	}
	cache := redis.NewClient(&redis.Options{
		Addr: f.config.CacheAddr,
	})
	handler := NewHandler(cache, r)

	rpc.Register(handler)
	rpc.Accept(inbound)
}
