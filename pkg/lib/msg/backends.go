package msg

type FeedBackend interface {
	Store(sessionID string, feedback *Feedback) error
}
