package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/jordan-wright/email"
	"github.com/labstack/gommon/log"
	"github.com/wneessen/go-mjml"
)

type emailService struct {
	templates  *template.Template
	fromEmail  string
	fromName   string
	asyncQueue AsyncQueue // Interface for async processing
}

func NewEmailService(templateDir string, fromEmail, fromName string, queue AsyncQueue) (EmailService, error) {
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.mjml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	return &emailService{
		templates:  templates,
		fromEmail:  fromEmail,
		fromName:   fromName,
		asyncQueue: queue,
	}, nil
}

func (s *emailService) SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error {
	// Create email task
	task := EmailTask{
		Subject:     subject,
		Template:    templatePath,
		ToEmail:     toEmail,
		Context:     context,
		Attachments: attachments,
		FromEmail:   s.fromEmail,
		FromName:    s.fromName,
	}

	// Queue the email task for async processing
	return s.asyncQueue.Enqueue("email", task)
}