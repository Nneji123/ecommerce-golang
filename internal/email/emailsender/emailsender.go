package emailsender

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// SendEmail sends an email using SMTP configuration and email details
func SendEmail(smtpConfig SMTPConfig, emailData Email) error {
	// Validate the From and To fields
	if emailData.From == "" {
		log.Println("From field is empty")
	}
	if len(emailData.To) == 0 {
		log.Println("At least one recipient must be specified in the To field")
	}
	// Initialize new email message
	msg := email.NewEmail()
	msg.From = emailData.From
	msg.To = emailData.To
	msg.Subject = emailData.Subject
	msg.Text = []byte(emailData.Body)
	msg.HTML = []byte(emailData.HTML)

	log.Println(emailData.BCC)
	msg.Cc = emailData.CC
	msg.Bcc = emailData.BCC

	// Add attachments
	for _, attachment := range emailData.Attachments {
		if _, err := msg.AttachFile(attachment); err != nil {
			return err
		}
	}

	// Connect to SMTP server
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Server)
	// Send email
	if err := msg.Send(smtpConfig.Server+":"+fmt.Sprint(smtpConfig.Port), auth); err != nil {
		return err
	}

	return nil
}
