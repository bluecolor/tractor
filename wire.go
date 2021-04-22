package tractor

type Wire interface {
	SendData(data *Data)
	SendFeed(feed *Feed)
}
