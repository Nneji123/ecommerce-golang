package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"mime/multipart"
	"net/smtp"
	"time"
)

type EmailMessage struct {
	from        string
	to          string
	subject     string
	htmlBody    string
	textBody    string
	attachments []Attachment
}

func NewEmailMessage() *EmailMessage {
	return &EmailMessage{}
}

func (m *EmailMessage) SetFrom(from string)              { m.from = from }
func (m *EmailMessage) SetTo(to string)                  { m.to = to }
func (m *EmailMessage) SetSubject(subject string)        { m.subject = subject }
func (m *EmailMessage) SetBody(html, text string)        { m.htmlBody = html; m.textBody = text }
func (m *EmailMessage) AddAttachments(atts []Attachment) { m.attachments = atts }

func (m *EmailMessage) Send(config *config.Config) error {
	tlsConfig := &tls.Config{
		ServerName:         config.SMTPHost,
		InsecureSkipVerify: true,
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort))
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer client.Close()

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(config.EmailFromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(m.to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %w", err)
	}
	defer writer.Close()

	var buf bytes.Buffer
	mimeWriter := multipart.NewWriter(&buf)

	headers := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n"+
			"Date: %s\r\n\r\n",
		m.from, m.to, m.subject,
		mimeWriter.Boundary(),
		time.Now().Format(time.RFC1123Z),
	)
	buf.WriteString(headers)

	textPart, err := mimeWriter.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}
	textPart.Write([]byte(m.textBody))

	htmlPart, err := mimeWriter.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}
	htmlPart.Write([]byte(m.htmlBody))

	for _, att := range m.attachments {
		attPart, err := mimeWriter.CreatePart(map[string][]string{
			"Content-Type":              {fmt.Sprintf("%s; name=%q", att.ContentType, att.Filename)},
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=%q", att.Filename)},
			"Content-Transfer-Encoding": {"base64"},
		})
		if err != nil {
			return fmt.Errorf("failed to create attachment part: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(att.Content)
		attPart.Write([]byte(encoded))
	}

	mimeWriter.Close()

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write email: %w", err)
	}

	return nil
}
