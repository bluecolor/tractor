package tractor

type Wire interface {
	SendFeed(feed Feed)
	SendData(data Data)
	ReadData() <-chan Data
	ReadFeeds() <-chan Feed
	CloseData()
	CloseFeed()
	Close()
}
