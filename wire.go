package tractor

type Wire interface {
	SendMessage(message *Message)
	ReadMessages() <-chan *Message
}
