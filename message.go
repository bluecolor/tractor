package tractor

type MessageType int
type FeedType int
type SenderType int
type Record []interface{}
type Data []Record

type ProgressFeed interface {
	Count() int
	Message() string
}
type progress struct {
	count   int
	message string
}

func (p progress) Count() int      { return p.count }
func (p progress) Message() string { return p.message }
func NewWriteProgress(count int, args ...string) Feed {
	var msg string = ""
	if len(args) > 0 {
		msg = args[0]
	}
	content := progress{
		count:   count,
		message: msg,
	}
	return Feed{
		Type:    Progress,
		Sender:  OutputPlugin,
		Content: content,
	}
}
func NewReadProgress(count int, args ...interface{}) Feed {
	var msg string

	if len(args) > 0 {
		msg = args[1].(string)
	}
	content := progress{
		count:   count,
		message: msg,
	}
	return Feed{
		Sender:  InputPlugin,
		Type:    Progress,
		Content: content,
	}
}

const (
	Success FeedType = iota
	Error
	Progress
)

const (
	Anonymous SenderType = iota
	InputPlugin
	OutputPlugin
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
