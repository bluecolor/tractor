package wire

import "github.com/bluecolor/tractor/lib/feed"

type Wire interface {
	SendFeed(f feed.Feed)
	SendData(d feed.Data)
	ReadData() <-chan feed.Data
	ReadFeeds() <-chan feed.Feed
	CloseData()
	CloseFeed()
	Close()
}
