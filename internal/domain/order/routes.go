package order

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
    // All order routes require authentication
    orders := e.Group("/orders")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	orders.Use(middleware.AuthMiddleware(cfg.JWTSecret))

    // User routes
    orders.POST("", h.Create)
    orders.GET("", h.ListUserOrders)
    orders.POST("/:id/cancel", h.CancelOrder)

    // Admin routes
    admin := orders.Group("", middleware.AdminMiddleware())
    admin.PUT("/:id/status", h.UpdateStatus)
}