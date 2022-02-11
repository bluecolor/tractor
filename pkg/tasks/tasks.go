package tasks

import (
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailSend     = "email:send"
	TypeExtractionRun = "extraction:run"
)

func NewExtractionRunTask(e models.Extraction) (*asynq.Task, error) {

	return nil, nil
	// payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	// if err != nil {
	// 	return nil, err
	// }
	// return asynq.NewTask(TypeEmailDelivery, payload), nil
}
