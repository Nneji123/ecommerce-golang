package user

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	auth := e.Group("/auth")

	// Public routes
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/confirm-registration", h.ConfirmRegistration)
	auth.POST("/password-reset-request", h.RequestPasswordReset)
	auth.POST("/confirm-password-reset", h.ConfirmPasswordReset)
}
