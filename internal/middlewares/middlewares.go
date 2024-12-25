package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nneji123/ecommerce-golang/config"
)

// corsWithConfig creates CORS middleware with custom configuration
func CorsWithConfig(config config.Config) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete},
	})
}

// rateLimiterMiddleware adds rate limiting to specific routes
func RateLimiterMiddleware(rateLimitedRoutes []string) echo.MiddlewareFunc {
	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 10, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	// Middleware function to apply rate limiting based on the route
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if the current route needs rate limiting
			for _, route := range rateLimitedRoutes {
				if route == c.Path() {
					return middleware.RateLimiterWithConfig(rateLimiterConfig)(next)(c)
				}
			}
			return next(c)
		}
	}
}
