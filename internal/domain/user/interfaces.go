package user

import (
	"github.com/nneji123/ecommerce-golang/internal/common/email"
	"github.com/nneji123/ecommerce-golang/internal/config"
)

type EmailService interface {
	SendVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
}

type emailService struct {
	service *email.EmailNotificationService
	config  *config.Config
}

func NewEmailService(service *email.EmailNotificationService, cfg *config.Config) EmailService { // Changed parameter type to pointer
	return &emailService{
		service: service,
		config:  cfg,
	}
}

func (s *emailService) SendVerificationEmail(email, token string) error {
	context := map[string]interface{}{
		"VerificationLink": s.config.AppURL + "/verify-email?token=" + token,
		"Token":            token,
	}

	return s.service.SendEmail(
		"Verify Your Email",
		"templates/verify-email.mjml",
		email,
		context,
		nil,
	)
}

func (s *emailService) SendPasswordResetEmail(email, token string) error {
	context := map[string]interface{}{
		"ResetLink": s.config.AppURL + "/reset-password?token=" + token,
		"Token":     token,
	}

	return s.service.SendEmail(
		"Reset Your Password",
		"templates/reset-password.mjml",
		email,
		context,
		nil,
	)
}
