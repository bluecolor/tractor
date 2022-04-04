package msg

type FeedBackend interface {
	Process(sessionID string, feedback *Feedback) error
}
