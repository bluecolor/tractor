package tasks

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/bluecolor/tractor/pkg/conf"
)

func TestServer(t *testing.T) {
	t.Parallel()
	mr := miniredis.RunT(t)
	c := conf.Tasks{
		Addr: mr.Addr(),
	}
	s := NewServer(c)

	go func(s *Server, t *testing.T) {
		if err := s.Run(); err != nil {
			t.Error(err)
		}
	}(s, t)
	time.Sleep(time.Second * 2)
	s.Shutdown()
}
