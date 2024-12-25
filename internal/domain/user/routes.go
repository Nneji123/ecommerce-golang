package user

import (
	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
	"log"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	// Public routes group
	auth := e.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/confirm-registration", h.ConfirmRegistration)
		auth.POST("/password-reset-request", h.RequestPasswordReset)
		auth.POST("/confirm-password-reset", h.ConfirmPasswordReset)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Define protected route group
	protected := e.Group("/user")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Register protected routes
	protected.GET("/detail", h.UserDetail)
}
