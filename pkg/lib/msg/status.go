package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type StatusFeed interface {
	Feed
	Status() types.Status
}

type statusFeed struct {
	feed
	status types.Status
}

func NewStatusFeed(s types.Sender, status types.Status, content ...interface{}) Feed {
	var c interface{}
	if len(content) > 0 {
		c = content[0]
	}
	return &statusFeed{
		feed:   feed{sender: s, content: c},
		status: status,
	}
}
func (f statusFeed) String() string {
	return fmt.Sprintf("%s: %s %s", f.sender, f.status, f.content)
}
func (f statusFeed) Status() types.Status {
	return f.status
}
