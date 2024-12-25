package middleware

import (
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

// Claims represents the JWT claims.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and extracts user claims
func AuthMiddleware(secretKey string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract the Authorization header
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
                return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
            }

            // Extract the token
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            // Parse and validate the token
            token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                // Validate the signing method
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return []byte(secretKey), nil
            })

            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }

            if !token.Valid {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }

            // Extract and validate claims
            claims, ok := token.Claims.(*Claims)
            if !ok {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
            }

            // Verify token expiration
            if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
                return echo.NewHTTPError(http.StatusUnauthorized, "token has expired")
            }

            // Store the claims in the context
            c.Set("user", claims)
            return next(c)
        }
    }
}

// Optional helper functions for role-based access control
func RequireRole(roles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*Claims)
            for _, role := range roles {
                if user.Role == role {
                    return next(c)
                }
            }
            return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
        }
    }
}

// Helper function to get user claims from context
func GetUserFromContext(c echo.Context) *Claims {
    user, ok := c.Get("user").(*Claims)
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
			userClaims, ok := c.Get("user").(*Claims)
			if !ok || userClaims.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}

			return next(c)
		}
	}
}

// GetUserClaims retrieves the user claims from the Echo context.
func GetUserClaims(c echo.Context) (*Claims, bool) {
	claims, ok := c.Get("user").(*Claims)
	return claims, ok
}
