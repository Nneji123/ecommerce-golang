package leads

import (
	"github.com/labstack/echo/v4"
)

func RegisterLeadHandlers(e *echo.Echo) {
	// Leads routes
	e.POST("/leads", HandlePostCreateLead)
	e.GET("/leads/:id", HandleGetLead)
	e.GET("/leads", HandleGetAllLeads)
	e.PUT("/leads/:id", HandleUpdateLead)
	e.DELETE("/leads/:id", HandleDeleteLead)

	e.POST("/leadlists", HandleCreateLeadList)
	e.GET("/leadlists/:id", HandleGetLeadList)
	e.PUT("/leadlists/:id", HandleUpdateLeadList)
	e.DELETE("/leadlists/:id", HandleDeleteLeadList)
}
