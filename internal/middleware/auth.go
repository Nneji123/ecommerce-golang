package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/domain/user"
)

// AuthMiddleware validates JWT tokens and extracts user claims.
func AuthMiddleware(authService user.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate the token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Store the claims in the Echo context
			c.Set("user", claims)
			return next(c)
		}
	}
}

// AdminOnly is a middleware that restricts access to admin users.
func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the user claims from the context
			userClaims, ok := c.Get("user").(*user.Claims)
			if !ok || userClaims.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}

			return next(c)
		}
	}
}

// GetUserClaims retrieves the user claims from the Echo context.
func GetUserClaims(c echo.Context) (*user.Claims, bool) {
	claims, ok := c.Get("user").(*user.Claims)
	return claims, ok
}
