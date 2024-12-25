package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/jordan-wright/email"
	"github.com/Boostport/mjml-go"
)



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
	// Render MJML template with provided context
	var mjmlBuf bytes.Buffer
	if err := w.templates.ExecuteTemplate(&mjmlBuf, task.Template, task.Context); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Convert MJML to HTML using mjml-go
	htmlContent, err := mjml.ToHTML(context.Background(), mjmlBuf.String(), mjml.WithMinify(true))
	if err != nil {
		var mjmlError mjml.Error
		if errors.As(err, &mjmlError) {
			return fmt.Errorf("MJML conversion error: %s - Details: %s", mjmlError.Message, mjmlError.Details)
		}
		return fmt.Errorf("failed to convert MJML to HTML: %w", err)
	}

	// Create and configure the email
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", task.FromName, task.FromEmail)
	e.To = []string{task.ToEmail}
	e.Subject = task.Subject
	e.HTML = []byte(htmlContent)

	// Add attachments, if any
	for _, att := range task.Attachments {
		if _, err := e.Attach(bytes.NewReader(att.Content), att.Filename, ""); err != nil {
			return fmt.Errorf("failed to attach file %s: %w", att.Filename, err)
		}
	}

	// Send email using SMTP
	err = e.Send(
		fmt.Sprintf("%s:%d", w.smtpHost, w.smtpPort),
		email.LoginAuth(w.smtpUser, w.smtpPass),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
