package leads

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/db"
)




// HandlePostCreateLead handles the endpoint for creating a new lead.
//
//	@Summary		Create a new lead
//	@Description	Creates a new lead based on the provided parameters.
//	@ID				handle-post-create-lead
//	@Accept			json
//	@Produce		json
//	@Param			body	body	db.Lead	true	"Lead details"
//	@Security		OAuth2Application[write] // Assuming OAuth2Application is the security scheme name and write is the required scope
//	@Success		201	{object}	db.Lead	"Created"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/leads [post]
func HandlePostCreateLead(c echo.Context) error {
	// Parse request body
	var lead db.Lead
	if err := c.Bind(&lead); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

    if err := db.CreateLead(&lead); err != nil {
        // Create an instance of ErrorResponse with the error message
        response := ErrorResponse{
            Error: "Failed to create lead",
        }
        // Return the JSON response
        return c.JSON(http.StatusInternalServerError, response)
    }


	return c.JSON(http.StatusCreated, lead)
}

//	@Summary		Retrieve a lead by ID
//	@Description	Retrieves a lead by its ID from the database.
//	@ID				handle-get-lead
//	@Param			id	path		int		true	"Lead ID"
//	@Success		200	{object}	db.Lead	"OK"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/leads/{id} [get]
func HandleGetLead(c echo.Context) error {
	// Parse lead ID from URL parameter
	leadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response:= ErrorResponse{Error: "Invalid lead ID"}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Retrieve lead from database
	lead, err := db.GetLead(uint(leadID))
	if err != nil {
	response:= ErrorResponse{Error: "Failed to retrieve lead"}

		return c.JSON(http.StatusInternalServerError, response)
	}
	if lead == nil {
		response:= ErrorResponse{Error: "Lead not found"}
		return c.JSON(http.StatusNotFound, response)
	}

	return c.JSON(http.StatusOK, lead)
}

//	@Summary		Retrieve all leads
//	@Description	Retrieves all leads from the database.
//	@ID				handle-get-all-leads
//	@Success		200	{object}	[]db.Lead	"OK"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/leads [get]
func HandleGetAllLeads(c echo.Context) error {
	// Retrieve all leads from the database
	leads, err := db.GetAllLeads()
	if err != nil {
		response:= ErrorResponse{Error: "Failed to retrieve leads"}

		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, leads)
}

//	@Summary		Update an existing lead
//	@Description	Updates an existing lead with the provided information.
//	@ID				handle-update-lead
//	@Param			id	path	int	true	"Lead ID"
//	@Accept			json
//	@Produce		json
//	@Param			body	body		db.Lead	true	"Updated lead details"
//	@Success		200		{object}	db.Lead	"OK"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/leads/{id} [put]
func HandleUpdateLead(c echo.Context) error {
	// Parse lead ID from URL parameter
	leadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response:= ErrorResponse{Error: "Invalid lead ID"}

		return c.JSON(http.StatusBadRequest, response)
	}

	// Retrieve lead from database
	lead, err := db.GetLead(uint(leadID))
	if err != nil {
		response:= ErrorResponse{Error: "Failed to retrieve lead"}

		return c.JSON(http.StatusInternalServerError, response)
	}
	if lead == nil {
		response:= ErrorResponse{Error: "Lead not found"}

		return c.JSON(http.StatusNotFound, response)
	}

	// Parse request body
	var updatedLead db.Lead
	if err := c.Bind(&updatedLead); err != nil {
		response:= ErrorResponse{Error: "Invalid request payload"}

		return c.JSON(http.StatusBadRequest, response)
	}

	// Update lead in database
	lead.Name = updatedLead.Name
	// Update other lead fields as needed...
	if err := db.UpdateLead(lead); err != nil {
		response:= ErrorResponse{Error: "Failed to update lead"}

		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.JSON(http.StatusOK, lead)
}

//	@Summary		Delete a lead by ID
//	@Description	Deletes a lead with the specified ID from the database.
//	@ID				handle-delete-lead
//	@Param			id	path	int	true	"Lead ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/leads/{id} [delete]
func HandleDeleteLead(c echo.Context) error {
	// Parse lead ID from URL parameter
	leadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response:= ErrorResponse{Error: "Invalid lead ID"}

		return c.JSON(http.StatusBadRequest, response)
	}

	// Delete lead from database
	if err := db.DeleteLead(uint(leadID)); err != nil {
		response:= ErrorResponse{Error: "Failed to delete lead"}

		return c.JSON(http.StatusInternalServerError, response)
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateLeadList creates a new lead list
func HandleCreateLeadList(c echo.Context) error {
	// Parse request body
	var leadList db.LeadList
	if err := c.Bind(&leadList); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Create lead list in database
	if err := db.CreateLeadList(&leadList); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create lead list"})
	}

	return c.JSON(http.StatusCreated, leadList)
}

// GetLeadList retrieves a lead list by ID
func HandleGetLeadList(c echo.Context) error {
	// Parse lead list ID from URL parameter
	leadListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lead list ID"})
	}

	// Retrieve lead list from database
	leadList, err := db.GetLeadListByID(uint(leadListID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve lead"})
	}
	if leadList == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lead list not found"})
	}

	return c.JSON(http.StatusOK, leadList)
}

// GetAllLeadListsHandler retrieves all lead lists
func HandleGetAllLeadLists(c echo.Context) error {
	// Retrieve all lead lists from the database
	leadLists, err := db.GetAllLeadLists()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve lead lists"})
	}

	return c.JSON(http.StatusOK, leadLists)
}

// UpdateLeadListHandler updates an existing lead list
func HandleUpdateLeadList(c echo.Context) error {
	// Parse lead list ID from URL parameter
	leadListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lead list ID"})
	}

	// Retrieve lead list from database
	leadList, err := db.GetLeadListByID(uint(leadListID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve lead list"})
	}
	if leadList == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Lead list not found"})
	}

	// Parse request body
	var updatedLeadList db.LeadList
	if err := c.Bind(&updatedLeadList); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Update lead list in database
	leadList.Name = updatedLeadList.Name
	// Update other lead list fields as needed...
	if err := db.UpdateLeadList(leadList); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update lead list"})
	}

	return c.JSON(http.StatusOK, leadList)
}

// DeleteLeadListHandler deletes a lead list by ID
func HandleDeleteLeadList(c echo.Context) error {
	// Parse lead list ID from URL parameter
	leadListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lead list ID"})
	}

	// Delete lead list from database
	if err := db.DeleteLeadList(uint(leadListID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete lead list"})
	}

	return c.NoContent(http.StatusNoContent)
}
