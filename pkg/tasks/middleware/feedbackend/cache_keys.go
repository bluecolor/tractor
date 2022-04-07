package feedbackend

import (
	"fmt"
)

func getSessionKey(sessionID string) string {
	return fmt.Sprintf("tractor:session:%s", sessionID)
}
func getPubsubKey(sessionID string) string {
	return fmt.Sprintf("tractor:session:%s:pubsub", sessionID)
}
