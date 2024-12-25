package emailpermutator

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GenerateEmailParams represents the parameters for generating email permutations
type GenerateEmailParams struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	NickName   string `json:"nickName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	Domain1    string `json:"domain1,omitempty"`
	Domain2    string `json:"domain2,omitempty"`
	Domain3    string `json:"domain3,omitempty"`
}

// HandlePostGenerateEmails generates email permutations based on the provided parameters.
//	@Summary		Generate email permutations
//	@Description	Generates email permutations based on the provided parameters.
//	@ID				handle-post-generate-emails
//	@Accept			json
//	@Produce		json
//	@Param			body	body	GenerateEmailParams	true	"Request body with parameters"
//	@Success		200		{array}	string				"OK"
//	@Router			/generate-emails [post]
func HandlePostGenerateEmails(c echo.Context) error {
	// Bind the request body to a struct
	var params GenerateEmailParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Generate email permutations
	emails := Permute(struct {
		FirstName  string
		LastName   string
		NickName   string
		MiddleName string
		Domain1    string
		Domain2    string
		Domain3    string
	}{
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		NickName:   params.NickName,
		MiddleName: params.MiddleName,
		Domain1:    params.Domain1,
		Domain2:    params.Domain2,
		Domain3:    params.Domain3,
	})

	// Return the email permutations as a JSON response
	return c.JSON(http.StatusOK, emails)
}
