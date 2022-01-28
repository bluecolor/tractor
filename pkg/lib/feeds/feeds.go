package feeds

import "fmt"

type (
	FeedType   int
	SenderType int
	Record     map[string]interface{}
	Data       []Record
	Progress   interface {
		Count() int
		Message() string
	}
)

const (
	SuccessFeed FeedType = iota
	ErrorFeed
	ProgressFeed
)

func (e FeedType) String() string {
	switch e {
	case SuccessFeed:
		return "SuccessFeed"
	case ErrorFeed:
		return "ErrorFeed"
	case ProgressFeed:
		return "ProgressFeed"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

const (
	SenderAnonymous SenderType = iota
	SenderInputConnector
	SenderOutputConnector
)

type Feed struct {
	Type    FeedType
	Sender  SenderType
	Content interface{}
}

func NewErrorFeed(sender SenderType, content interface{}) Feed {
	return Feed{
		Type:    ErrorFeed,
		Sender:  sender,
		Content: content,
	}
}
func NewSuccessFeed(sender SenderType) Feed {
	return Feed{
		Type:   SuccessFeed,
		Sender: sender,
	}
}
func NewFeed(sender SenderType, feedType FeedType, content interface{}) Feed {
	return Feed{
		Type:    feedType,
		Sender:  sender,
		Content: content,
	}
}

type ProgressMessage struct {
	count   int
	message string
}

func (p *ProgressMessage) Count() int {
	return p.count
}
func (p *ProgressMessage) Message() string {
	return p.message
}
func NewWriteProgress(count int, args ...interface{}) Feed {
	var msg string

	if len(args) > 0 {
		msg = args[1].(string)
	}
	content := ProgressMessage{
		count:   count,
		message: msg,
	}
	return Feed{
		Sender:  SenderOutputConnector,
		Type:    ProgressFeed,
		Content: content,
	}
}
func NewReadProgress(count int, args ...interface{}) Feed {
	var msg string

	if len(args) > 0 {
		msg = args[1].(string)
	}
	content := ProgressMessage{
		count:   count,
		message: msg,
	}
	return Feed{
		Sender:  SenderInputConnector,
		Type:    ProgressFeed,
		Content: content,
	}
}
