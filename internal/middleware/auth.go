package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/nneji123/ecommerce-golang/internal/common/models"
)

// AuthMiddleware validates JWT tokens and extracts user claims.
func AuthMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the Authorization header.
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
			}

			// Extract the token.
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate the token.
			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// Verify token expiration (optional since jwt.ParseWithClaims handles it).
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
			}

			// Set claims in the context for later use.
			c.Set("userClaims", claims)

			// Proceed to the next handler.
			return next(c)
		}
	}
}

// Helper function to retrieve user claims from the context.
func GetUserFromContext(c echo.Context) *models.Claims {
	user, ok := c.Get("user").(*models.Claims)
	if !ok {
		return nil
	}
	return user
}

// AdminOnly is a middleware that restricts access to admin users.
func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the user claims from the context
			userClaims, ok := c.Get("user").(*models.Claims)
			if !ok || userClaims.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}

			return next(c)
		}
	}
}

// GetUserClaims retrieves the user claims from the Echo context.
func GetUserClaims(c echo.Context) (*models.Claims, bool) {
	claims, ok := c.Get("user").(*models.Claims)
	return claims, ok
}
