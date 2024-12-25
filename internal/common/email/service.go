package email

import (
	"fmt"
	"github.com/nneji123/ecommerce-golang/internal/config"
)

type EmailNotificationService struct {
	service EmailService
}

func NewEmailNotificationService(serviceType string, config *config.Config) (*EmailNotificationService, error) {
	var service EmailService

	switch serviceType {
	case "smtp":
		service = NewSMTPService(config)
	case "sendgrid":
		service = NewSendGridService(config)
	default:
		return nil, fmt.Errorf("unsupported email service type: %s", serviceType)
	}

	return &EmailNotificationService{service: service}, nil
}

func (s *EmailNotificationService) SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error {
	return s.service.SendEmail(subject, templatePath, toEmail, context, attachments)
}

// SMTP Service implementation
type SMTPService struct {
	config *config.Config
}

func NewSMTPService(config *config.Config) *SMTPService {
	return &SMTPService{config: config}
}

func (s *SMTPService) SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error {
	html, text, err := RenderTemplate(templatePath, context)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	msg := NewEmailMessage()
	msg.SetFrom(fmt.Sprintf("%s <%s>", s.config.EmailFromName, s.config.EmailFromAddress))
	msg.SetTo(toEmail)
	msg.SetSubject(subject)
	msg.SetBody(html, text)
	msg.AddAttachments(attachments)

	if err := msg.Send(s.config); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendGrid Service implementation
type SendGridService struct {
	config *config.Config
}

func NewSendGridService(config *config.Config) *SendGridService {
	return &SendGridService{config: config}
}

func (s *SendGridService) SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error {
	// Implementation similar to SMTPService but using SendGrid API
	return nil
}
