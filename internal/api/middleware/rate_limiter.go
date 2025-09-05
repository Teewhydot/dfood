package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/httprate"
)

// RateLimitMiddleware creates a Gin middleware function for rate limiting using httprate
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	// Create httprate limiter with IP-based limiting
	limiter := httprate.Limit(
		limit,    // number of requests
		window,   // time window
		httprate.WithKeyFuncs(httprate.KeyByIP), // rate limit by IP
	)
	
	return func(c *gin.Context) {
		// Create a dummy handler that continues to the next middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})
		
		// Apply the rate limiter
		limitedHandler := limiter(handler)
		limitedHandler.ServeHTTP(c.Writer, c.Request)
	}
}
