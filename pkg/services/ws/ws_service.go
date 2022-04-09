package ws

import (
	"net/http"

	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
)

type Service struct {
	client *redis.Client
}

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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	ps := s.client.Subscribe("tractor:session:feeds")
	defer func() {
		_ = ps.Close()
		_ = conn.Close()
	}()

	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
		if msg == nil || msg.Payload == "" || msg.Payload == "null" {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			utils.ErrorWithJSON(w, http.StatusInternalServerError, err)
			return
		}
	}
}
