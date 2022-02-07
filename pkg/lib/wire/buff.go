package wire

import (
	"sync"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

const BufferSize = 1000 // todo from .env

type BufferedWire struct {
	mu sync.Mutex
	*Wire
	buffer []msg.Record
	size   int
}

func NewBuffered(w *Wire, size ...int) BufferedWire {
	var sz int = BufferSize
	if len(size) > 0 {
		sz = size[0]
	}
	return BufferedWire{
		Wire:   w,
		buffer: []msg.Record{},
		size:   sz,
	}
}
func (bw *BufferedWire) Send(data interface{}) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	var box []msg.Record
	switch val := data.(type) {
	case []msg.Record:
		box = val
	case msg.Record:
		box = []msg.Record{val}
	default:
		return
	}
	for _, r := range box {
		bw.buffer = append(bw.buffer, r)
		if len(bw.buffer) >= bw.size {
			bw.Wire.SendData(bw.buffer)
			bw.buffer = []msg.Record{}
		}
	}
}
func (bw *BufferedWire) Flush() {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if len(bw.buffer) > 0 {
		bw.Wire.SendData(bw.buffer)
		bw.buffer = []msg.Record{}
	}
}
