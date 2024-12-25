package email

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/wneessen/go-mjml"
)

// NewEmailService initializes the email service
func NewEmailService(templateDir string, fromEmail, fromName string, queue AsyncQueue) (*emailService, error) {
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

// SendEmail queues an email for asynchronous processing
func (s *emailService) SendEmail(subject, templatePath string, toEmail string, context map[string]interface{}, attachments []Attachment) error {
	// Create the email task
	task := EmailTask{
		Subject:     subject,
		Template:    templatePath,
		ToEmail:     toEmail,
		Context:     context,
		Attachments: attachments,
		FromEmail:   s.fromEmail,
		FromName:    s.fromName,
	}

	// Queue the email task
	return s.asyncQueue.Enqueue("email:send", task)
}

// EnqueueTask handles queuing an email task with asynq
func (q *asynqQueue) Enqueue(queue string, task EmailTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal email task: %w", err)
	}

	// Create an Asynq task with the payload
	asynqTask := asynq.NewTask(queue, payload)

	// Add the task to the queue
	_, err = q.client.Enqueue(asynqTask)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// asynqQueue wraps the Asynq client for task enqueuing
type asynqQueue struct {
	client *asynq.Client
}

// NewAsynqQueue creates a new instance of asynqQueue
func NewAsynqQueue(client *asynq.Client) *asynqQueue {
	return &asynqQueue{client: client}
}
