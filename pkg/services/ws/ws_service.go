package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Service struct {
	client *redis.Client
}

var ctx = context.Background()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewService(client *redis.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) SubSessionFeeds(w http.ResponseWriter, r *http.Request) {
	// todo CheckOrigin
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	ps := s.client.Subscribe(ctx, "tractor:session:feeds")
	defer func() {
		_ = ps.Close()
		_ = conn.Close()
	}()

	for {
		m, err := ps.ReceiveMessage(ctx)
		if err != nil {
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		if m == nil || m.Payload == "" || m.Payload == "null" {
			continue
		}
		feed := &msg.Feed{}
		if err := feed.Unmarshal([]byte(m.Payload)); err != nil {
			log.Error().Err(err).Msg("unmarshal feed")
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		payload, err := json.Marshal(map[string]interface{}{
			"_":         feed,
			"sessionId": feed.SessionID,
			"type":      feed.Type.String(),
			"sender":    feed.Sender.String(),
			"content":   feed.Content,
		})
		if err != nil {
			log.Error().Err(err).Msg("marshal feed")
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
	}
}
