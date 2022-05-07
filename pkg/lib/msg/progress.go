package msg

import (
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
)

type ProgressFeed interface {
	Feed
	Progress() int
}
type progressFeed struct {
	feed
}

func NewProgressFeed(sender types.Sender, progress int) Feed {
	return &progressFeed{
		feed: feed{sender: sender, content: progress},
	}
}
func (f progressFeed) Progress() int {
	return f.content.(int)
}
func (f progressFeed) String() string {
	return fmt.Sprintf("progress - %s: %d", f.sender, f.content)
}
