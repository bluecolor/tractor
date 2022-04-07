package feedbackend

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/go-redis/redis/v7"
)

type FeedBackend struct {
	cache *redis.Client
	repo  *repo.Repository
}

func New(cacheAddr string, dbConfig conf.DB) (*FeedBackend, error) {
	r, err := repo.New(dbConfig)
	if err != nil {
		return nil, err
	}
	return &FeedBackend{
		cache: redis.NewClient(&redis.Options{
			Addr: cacheAddr,
		}),
		repo: r,
	}, nil
}
func (f *FeedBackend) CloseCache() error {
	return f.cache.Close()
}

func (f *FeedBackend) Process(sessionID string, feed *msg.Feed) error {
	if err := f.UpdateRepo(sessionID, feed); err != nil {
		return err
	}
	return nil
}
