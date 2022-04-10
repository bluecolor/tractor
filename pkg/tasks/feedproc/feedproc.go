package feedproc

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type FeedProcessor struct {
	cache *redis.Client
	repo  *repo.Repository
}

func New(config conf.FeedProcessor) (*FeedProcessor, error) {
	r, err := repo.New(config.DB)
	if err != nil {
		log.Error().Err(err).Msg("failed to create repo")
		return nil, err
	}
	cache := redis.NewClient(&redis.Options{
		Addr: config.CacheAddr,
	})
	return &FeedProcessor{
		cache: cache,
		repo:  r,
	}, nil
}

func (fp *FeedProcessor) Process(feed *msg.Feed) error {
	if err := fp.UpdateCache(feed); err != nil {
		return err
	}
	if err := fp.UpdateRepo(feed); err != nil {
		return err
	}
	if err := fp.Publish(feed); err != nil {
		return err
	}
	return nil
}
