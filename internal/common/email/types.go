package email

import "html/template"

type Attachment struct {
	Filename string
	Content  []byte
}

type EmailService interface {
	SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error
}

// EmailTask represents the structure of an email task
type EmailTask struct {
	Subject     string                 `json:"subject"`
	Template    string                 `json:"template"`
	ToEmail     string                 `json:"to_email"`
	Context     map[string]interface{} `json:"context"`
	Attachments []Attachment           `json:"attachments"`
	FromEmail   string                 `json:"from_email"`
	FromName    string                 `json:"from_name"`
}

type EmailWorker struct {
	templates *template.Template
	smtpHost  string
	smtpPort  int
	smtpUser  string
	smtpPass  string
}

type AsyncQueue interface {
	Enqueue(queue string, task EmailTask) error
}

type emailService struct {
	templates  *template.Template
	fromEmail  string
	fromName   string
	asyncQueue AsyncQueue
}