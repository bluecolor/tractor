package tractor

type Wire interface {
	SendFeed(sender SenderType, feed *Feed)
	SendData(data *Data)
	ReadData() <-chan *Data
	ReadFeeds() <-chan *Feed
	Close()
}
