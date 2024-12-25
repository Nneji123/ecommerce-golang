package email

type Attachment struct {
	Filename    string
	Content     []byte
	ContentType string
}

type EmailService interface {
	SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error
}
