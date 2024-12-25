package emailsender

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

// NewHandler creates a new instance of Handler
func NewHandler(client *asynq.Client) *Handler {
	return &Handler{Client: client}
}

//

//	@Summary		Send emails
//	@Description	Parses the multipart form data into an EmailRequest struct and enqueues an email delivery task.
//	@ID				handle-post-send-emails
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			from		formData	string		true	"From email address"
//	@Param			to			formData	[]string	true	"List of recipient email addresses"
//	@Param			cc			formData	[]string	false	"List of CC email addresses"
//	@Param			bcc			formData	[]string	false	"List of BCC email addresses"
//	@Param			subject		formData	string		true	"Email subject"
//	@Param			body		formData	string		true	"Email body"
//	@Param			html		formData	string		true	"HTML Email"
//	@Param			attachments	formData	[]file		false	"List of file attachments"
//	@Param			server		formData	string		true	"SMTP server address"
//	@Param			port		formData	string		true	"SMTP server port"
//	@Param			username	formData	string		true	"SMTP username"
//	@Param			password	formData	string		true	"SMTP password"
//	@Success		200			{object}	JSONResponse
//	@Router			/send-email [post]
//
// HandlePostSendEmails handles the endpoint for sending emails.
func (h *Handler) HandlePostSendEmails(c echo.Context) error {
	emailRequest := new(EmailRequest)
	if err := c.Bind(emailRequest); err != nil {
		return err
	}

	// Get the value of "bcc" from the form
	bccValue := c.FormValue("bcc")
	ccValue := c.FormValue("cc")

	// Split the string into individual email addresses
	bccList := strings.Split(bccValue, ",")
	ccList := strings.Split(ccValue, ",")

	// Assign the slice of email addresses to EmailRequest's BCC field
	emailRequest.BCC = bccList
	emailRequest.CC = ccList

	// Parse SMTP configuration fields
	portStr := c.FormValue("smtpConfig.port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("Error converting to integer")
		return err
	}

	smtpConfig := SMTPConfig{
		Server:   c.FormValue("smtpConfig.server"),
		Port:     port,
		Username: c.FormValue("smtpConfig.username"),
		Password: c.FormValue("smtpConfig.password"),
	}

	// Assign SMTP configuration to EmailRequest struct
	emailRequest.SMTPConfig = smtpConfig

	// Check if files are present in the request
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["attachments"]
	var attachments []string

	destinationFolder := "./data"

	// Iterate through each file
	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		destinationPath := filepath.Join(destinationFolder, file.Filename)
		dst, err := os.Create(destinationPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Append the destination path to the attachments slice
		attachments = append(attachments, destinationPath)
	}

	// Update EmailRequest with attachments
	emailRequest.Attachments = attachments

	// Marshal email data payload
	payload, err := json.Marshal(emailRequest)

	if err != nil {
		response := ErrorResponse{
			Error: fmt.Sprintf("Error occured with json encoding: %v\n", err),
		}
		c.JSON(http.StatusBadRequest, response)
		log.Printf("Error occured with json marshal %v", err)
	}

	// Create a new email delivery task
	task := asynq.NewTask(TypeEmailDelivery, payload)
	// Enqueue the task
	var response interface{}

	if info, err := h.Client.Enqueue(task); err != nil {
		response = ErrorResponse{
			Error: fmt.Sprintf("Error enqueuing email delivery task: %v\n", err),
		}
		c.JSON(http.StatusInternalServerError, response)
		log.Printf("Error enqueuing email delivery task: %v\n", err)
		return err
	} else {
		response = JSONResponse{
			Message: fmt.Sprintf("Successfully enqueued email sending task: ID=%s, Queue=%s", info.ID, info.Queue),
		}
		log.Printf("Enqueued task: ID=%s, Queue=%s", info.ID, info.Queue)
	}

	// Return the JSON response
	return c.JSON(http.StatusOK, response)
}
