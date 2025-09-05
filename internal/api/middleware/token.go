package middleware

import (
	"dfood/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(context *gin.Context) {
	// Implement token authentication logic here
	token := context.GetHeader("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}
	// Validate the token
	if _, err := utils.ValidateToken(token); err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}
	context.Next()
}
