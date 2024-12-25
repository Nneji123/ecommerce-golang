package product

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/config"
	"github.com/nneji123/ecommerce-golang/internal/middleware"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	products := e.Group("/products")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	products.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Public routes (authenticated users)
	products.GET("", h.List)
	products.GET("/:id", h.Get)

	// Admin routes
	admin := products.Group("", middleware.AdminMiddleware())
	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}
