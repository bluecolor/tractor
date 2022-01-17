package feed

type (
	FeedType     int
	SenderType   int
	Record       []interface{}
	Data         []Record
	ProgressFeed interface {
		Count() int
		Message() string
	}
)

const (
	Success FeedType = iota
	Error
	Progress
)
const (
	SenderAnonymous SenderType = iota
	SenderInputPlugin
	SenderOutputPlugin
)

type Feed struct {
	Type    FeedType
	Sender  SenderType
	Content interface{}
}

func NewErrorFeed(sender SenderType, content interface{}) Feed {
	return Feed{
		Type:    Error,
		Sender:  sender,
		Content: content,
	}
}
func NewSuccessFeed(sender SenderType) Feed {
	return Feed{
		Type:   Success,
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
		Sender:  SenderOutputPlugin,
		Type:    Progress,
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
		Sender:  SenderInputPlugin,
		Type:    Progress,
		Content: content,
	}
}
