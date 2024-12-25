package emailvalidator

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type validateRequest struct {
	Email      string `query:"email"`
	SMTPCheck  bool   `query:"smtp,omitempty"`
	SOCKSCheck bool   `query:"socks,omitempty"`
}

type validateResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// HandleGetValidateEmail handles the endpoint for validating emails.
//	@Summary		Validate emails
//	@Description	Validates email addresses based on the provided parameters.
//	@ID				handle-get-validate-email
//	@Accept			json
//	@Produce		json
//	@Param			email	query	string	true	"Email address to validate"
//	@Param			smtp	query	boolean	false	"Perform SMTP check"	default(false)
//	@Param			socks	query	boolean	false	"Perform SOCKS check"	default(false)
//	@Security		OAuth2Application[callback]
//	@Success		200	{object}	validateResponse	"OK"
//	@Failure		400	{object}	validateResponse	"Bad Request"
//	@Failure		500	{object}	validateResponse	"Internal Server Error"
//	@Router			/validate-email [get]
func HandleGetValidateEmail(c echo.Context) error {
	req := new(validateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, validateResponse{Error: err.Error()})
	}

	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, validateResponse{Error: "email parameter is required"})
	}

	var result interface{}
	var err error

	if req.SMTPCheck {
		result, err = VerifyWithSMTPCheck("", "")
	} else if req.SOCKSCheck {
		result, err = VerifyWithSOCKS("", "")
	} else {
		result, err = VerifyEmail(req.Email)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, validateResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, validateResponse{Result: result})
}
