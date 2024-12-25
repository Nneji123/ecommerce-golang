package main

import (
	"fmt"
	"github.com/nneji123/ecommerce-golang/internal/common/email"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the EmailNotificationService with "smtp" service type
	emailService, err := email.NewEmailNotificationService("smtp", &cfg)
	if err != nil {
		log.Fatalf("Error initializing email service: %v", err)
	}

	// Define email context data (you can add any dynamic data for your templates)
	context := map[string]interface{}{
		"FirstName": "John",
		"LastName":  "Doe",
	}

	// Send an email
	subject := "Welcome to Our Service!"
	templatePath := "templates/welcome.mjml" // Ensure this is the correct path
	toEmail := "recipient@example.com"

	err = emailService.SendEmail(subject, templatePath, toEmail, context, nil)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
