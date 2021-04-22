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
