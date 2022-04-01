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
	c := conf.Worker{
		Addr: mr.Addr(),
	}
	w := NewWorker(c)
	go func(s *Worker, t *testing.T) {
		if err := s.Start(); err != nil {
			t.Error(err)
		}
	}(w, t)
	time.Sleep(time.Second * 2)
	w.Shutdown()
}
