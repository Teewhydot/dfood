package middleware

import (
	"net/http"
	"strings"

	"dfood/internal/repository"
	"dfood/internal/utils"
	"dfood/pkg/errors"
	"dfood/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
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

		// Fetch user from database to get user ID
		user, err := userRepo.GetByEmail(email)
		if err != nil {
			result := errors.HandleError(
				func() (interface{}, error) {
					return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not found", err)
				},
				"fetching user by email",
			)
			result.RespondWithJSON(c)
			c.Abort()
			return
		}

		// Set user context for use in handlers
		c.Set("user_id", user.ID)
		c.Set("user_email", email)
		c.Set("token", token)
		logger.Info("User ID: ", user.ID)
		logger.Info("User email: ", email)
		logger.Info("User access token: ", token)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func OptionalAuthMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
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
			// Try to fetch user from database to get user ID
			user, err := userRepo.GetByEmail(email)
			if err == nil {
				// Set user context for use in handlers
				c.Set("user_id", user.ID)
				c.Set("user_email", email)
				c.Set("token", token)
			}
		}

		c.Next()
	}
}
