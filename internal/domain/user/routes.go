package user

import (
	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
	"github.com/nneji123/ecommerce-golang/internal/config"
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

    // Protected routes group
    user := e.Group("/user")
    user.Use(middleware.AuthMiddleware(cfg.JWTSecret))
    {
        user.GET("/detail", h.UserDetail)
        
        // Admin-only routes (if you have any)
        admin := user.Group("")
        admin.Use(middleware.RequireRole("admin"))
        {
            // Add admin routes here
            // admin.GET("/all-users", h.ListAllUsers)
        }
    }
}