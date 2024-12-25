package service
// internal/shared/email/types.go
// package email

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
)

type Attachment struct {
	Filename string
	Content  []byte
}

type EmailService interface {
	SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error
}

// internal/shared/email/service.go
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

// internal/shared/email/worker.go
package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/jordan-wright/email"
	"github.com/wneessen/go-mjml"
)

type EmailTask struct {
	Subject     string
	Template    string
	ToEmail     string
	Context     map[string]interface{}
	Attachments []Attachment
	FromEmail   string
	FromName    string
}

type EmailWorker struct {
	templates *template.Template
	smtpHost  string
	smtpPort  int
	smtpUser  string
	smtpPass  string
}

func NewEmailWorker(templateDir, smtpHost string, smtpPort int, smtpUser, smtpPass string) (*EmailWorker, error) {
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.mjml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email templates: %w", err)
	}

	return &EmailWorker{
		templates: templates,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
		smtpUser:  smtpUser,
		smtpPass:  smtpPass,
	}, nil
}

func (w *EmailWorker) ProcessEmail(task EmailTask) error {
	// Render MJML template
	var mjmlBuf bytes.Buffer
	if err := w.templates.ExecuteTemplate(&mjmlBuf, task.Template, task.Context); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Convert MJML to HTML
	htmlContent, err := mjml.New().Convert(mjmlBuf.String())
	if err != nil {
		return fmt.Errorf("failed to convert MJML to HTML: %w", err)
	}

	// Create email
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", task.FromName, task.FromEmail)
	e.To = []string{task.ToEmail}
	e.Subject = task.Subject
	e.HTML = []byte(htmlContent)

	// Add attachments
	for _, att := range task.Attachments {
		_, err := e.Attach(bytes.NewReader(att.Content), att.Filename, "")
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", att.Filename, err)
		}
	}

	// Send email
	return e.Send(
		fmt.Sprintf("%s:%d", w.smtpHost, w.smtpPort),
		email.LoginAuth(w.smtpUser, w.smtpPass),
	)
}