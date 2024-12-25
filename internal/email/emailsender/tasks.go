package emailsender

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	// Unmarshal payload into EmailRequest
	var emailRequest EmailRequest
	if err := json.Unmarshal(t.Payload(), &emailRequest); err != nil {
		log.Printf("failed to unmarshal email payload: %v", err)
		return err
	}

	// Send email using SendEmail function
	email := Email{
		From:        emailRequest.From,
		To:          emailRequest.To,
		CC:          emailRequest.CC,
		BCC:         emailRequest.BCC,
		Subject:     emailRequest.Subject,
		Body:        emailRequest.Body,
		HTML:        emailRequest.HTML,
		Attachments: emailRequest.Attachments,
	}
	if err := SendEmail(emailRequest.SMTPConfig, email); err != nil {
		log.Printf("failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully: To=%v, Subject=%s", email.To, email.Subject)
	return nil
}
