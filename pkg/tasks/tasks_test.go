package tasks

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/bluecolor/tractor/pkg/conf"
)

func TestWorker(t *testing.T) {
	t.Parallel()
	mr := miniredis.RunT(t)
	c := conf.Tasks{
		Addr: mr.Addr(),
	}
	w := NewWorker(c)
	go func(s *Server, t *testing.T) {
		if err := s.Run(); err != nil {
			t.Error(err)
		}
	}(w, t)
	time.Sleep(time.Second * 2)
	w.Shutdown()
}
