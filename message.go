package tractor

type MessageType int
type FeedType int
type SenderType int
type Record []interface{}
type Data []Record

type ProgressFeed interface {
	Count() int
	Total() int
	Message() string
}
type progress struct {
	count   int
	total   int
	message string
}

func (p *progress) Count() int      { return p.count }
func (p *progress) Total() int      { return p.total }
func (p *progress) Message() string { return p.message }
func NewWriteProgress(count int, args ...string) *Message {
	var msg string = ""
	if len(args) > 0 {
		msg = args[0]
	}
	content := progress{
		count:   count,
		message: msg,
	}
	feed := Feed{
		Type:    Progress,
		Content: content,
	}
	return &Message{
		Sender:  OutputPlugin,
		Type:    FeedMessage,
		Content: feed,
	}
}
func NewReadProgress(count int, args ...interface{}) *Message {
	var msg string
	var total int

	if len(args) > 0 {
		total = args[0].(int)
		if len(args) > 1 {
			msg = args[1].(string)
		}
	}
	content := progress{
		count:   count,
		total:   total,
		message: msg,
	}
	feed := Feed{
		Type:    Progress,
		Content: content,
	}
	return &Message{
		Sender:  InputPlugin,
		Type:    FeedMessage,
		Content: feed,
	}
}

const (
	Success FeedType = iota
	Error
	Progress
)

const (
	FeedMessage MessageType = iota
	DataMessage
	CatalogMessage
)

const (
	Anonymous SenderType = iota
	InputPlugin
	OutputPlugin
)

type Feed struct {
	Type    FeedType
	Content interface{}
}

type Message struct {
	Type    MessageType
	Sender  SenderType
	Content interface{}
}

func NewErrorFeed(sender SenderType, content interface{}) *Message {

	feed := Feed{
		Type:    Error,
		Content: content,
	}
	return &Message{
		Type:    FeedMessage,
		Sender:  sender,
		Content: feed,
	}
}

func NewDataMessage(data Data) *Message {
	return &Message{
		Type:    FeedMessage,
		Sender:  InputPlugin,
		Content: data,
	}
}

func NewSuccessFeed(sender SenderType) *Message {
	feed := Feed{
		Type: Success,
	}
	return &Message{
		Type:    FeedMessage,
		Sender:  sender,
		Content: feed,
	}
}

func NewCatalogMessage(sender SenderType) {}
