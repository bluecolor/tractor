package tractor

type MessageType int
type FeedType int
type SenderType int
type Record []interface{}
type Data []Record

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
