package routes

import (
	"time"

	"dfood/internal/api/handlers"
	"dfood/internal/api/middleware"
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthService  service.AuthService
}

func SetupRoutes(deps *Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Apply global rate limiting: per-IP, per-endpoint using httprate
	router.Use(middleware.RateLimitMiddleware(10, time.Minute)) // 10 requests per minute per IP

	authHandler := handlers.NewAuthHandler(deps.AuthService)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			// update user password endpoint
			auth.POST("/password/update", authHandler.UpdatePassword)
		}

	}

	return router
}
