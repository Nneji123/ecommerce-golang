package emailsender

import (
	"github.com/hibiken/asynq"
)

// Email represents the details of an email
type Email struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	CC          []string `json:"cc,omitempty"`
	BCC         []string `json:"bcc,omitempty"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	HTML        string   `json:"html,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

// SMTPConfig represents SMTP configuration details
type SMTPConfig struct {
	Server   string `form:"server"`
	Port     int    `form:"port"`
	Username string `form:"username"`
	Password string `form:"password"`
}

// Handler represents the handler for email sending functionality
type Handler struct {
	Client *asynq.Client
}

type EmailRequest struct {
	From        string     `form:"from"`
	To          []string   `form:"to"`
	CC          []string   `form:"cc,omitempty"`
	BCC         []string   `form:"bcc,omitempty"`
	Subject     string     `form:"subject"`
	Body        string     `form:"body"`
	HTML        string     `form:"html"`
	Attachments []string   `form:"attachments,omitempty"`
	SMTPConfig  SMTPConfig `form:"smtpConfig"`
}

type JSONResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	TypeEmailDelivery = "email:deliver"
)
