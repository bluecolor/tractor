package feedbackend

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/repo"
	"github.com/go-redis/redis/v7"
)

func getSessionKey(sessionID string) string {
	return fmt.Sprintf("tractor:session:%s", sessionID)
}
func getPubsubKey() string {
	return fmt.Sprintf("tractor:session:feeds")
}

type Handler struct {
	cache *redis.Client
	repo  *repo.Repository
}

func NewHandler(cache *redis.Client, repo *repo.Repository) *Handler {
	return &Handler{
		cache: cache,
		repo:  repo,
	}
}

func (h *Handler) Process(feed *msg.Feed, args ...interface{}) error {
	if err := h.UpdateCache(feed); err != nil {
		return err
	}
	if err := h.UpdateRepo(feed); err != nil {
		return err
	}
	return nil
}
