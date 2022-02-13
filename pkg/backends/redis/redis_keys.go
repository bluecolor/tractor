package redis

import "fmt"

const (
	StatusKey         = "status"
	InputProgressKey  = "input_progress"
	OutputProgressKey = "output_progress"
	ErrorKey          = "error"
	WarningKey        = "warning"
	PubsubKey         = "pubsub"
)

func getSessionKey(sessionID string) string {
	return fmt.Sprintf("tractor:session:%s", sessionID)
}
func getPubsubKey(sessionID string) string {
	return fmt.Sprintf("tractor:session:%s:%s", sessionID, PubsubKey)
}
