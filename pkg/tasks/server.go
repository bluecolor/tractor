package tasks

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

type Server struct {
	config *conf.Tasks
	server *asynq.Server
}

func NewServer(c conf.Tasks) *Server {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Addr},
		asynq.Config{
			Concurrency: c.Concurrency,
		},
	)
	return &Server{
		config: &c,
		server: server,
	}
}

func (s *Server) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeExtractionRun, HandleExtractionTask)
	if err := s.server.Run(mux); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}
