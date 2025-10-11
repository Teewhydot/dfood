package middleware

import (
	"net/http"
	"strings"

	"dfood/internal/utils"
	"dfood/pkg/errors"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			result := errors.HandleError(
				func() (interface{}, error) {
					return nil, errors.NewHTTPError(http.StatusUnauthorized, "Authorization header required", nil)
				},
				"missing authorization header",
			)
			result.RespondWithJSON(c)
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			result := errors.HandleError(
				func() (interface{}, error) {
					return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token", err)
				},
				"validating JWT token",
			)
			result.RespondWithJSON(c)
			c.Abort()
			return
		}

		// Extract user email from token
		email, ok := (*claims)["sub"].(string)
		if !ok {
			result := errors.HandleError(
				func() (interface{}, error) {
					return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid token claims", nil)
				},
				"extracting user email from token",
			)
			result.RespondWithJSON(c)
			c.Abort()
			return
		}

		// Set user context for use in handlers
		c.Set("user_email", email)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			// Token is invalid, but we continue without setting user context
			c.Next()
			return
		}

		// Extract user email from token
		email, ok := (*claims)["sub"].(string)
		if ok {
			// Set user context for use in handlers
			c.Set("user_email", email)
			c.Set("token", token)
		}

		c.Next()
	}
}
