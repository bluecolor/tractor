package msg

type FeedbackBackend interface {
	Store(sessionID string, feedback *Feedback) error
}
