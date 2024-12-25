package email

import (
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/email/emailpermutator"
	"github.com/nneji123/ecommerce-golang/internal/email/emailscraper"
	"github.com/nneji123/ecommerce-golang/internal/email/emailsender"
	"github.com/nneji123/ecommerce-golang/internal/email/emailvalidator"
)

func RegisterEmailHandlers(e *echo.Echo, client *asynq.Client) {
	emailHandler := emailsender.NewHandler(client)
	e.POST("/generate-emails", emailpermutator.HandlePostGenerateEmails)
	e.POST("/send-email", emailHandler.HandlePostSendEmails)
	e.POST("/scrape-emails", emailscraper.ScrapeEmails)
	e.GET("/validate-email", emailvalidator.HandleGetValidateEmail)
}
