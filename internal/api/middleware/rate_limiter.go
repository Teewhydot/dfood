package middleware

import (
	"net/http"
	"strconv"
	"time"

	"dfood/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/httprate"
)

// RateLimitMiddleware creates a Gin middleware function for rate limiting using httprate
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	// Create httprate limiter with IP-based limiting and custom options
	limiter := httprate.Limit(
		limit,                                   // number of requests
		window,                                  // time window
		httprate.WithKeyFuncs(httprate.KeyByIP), // rate limit by IP
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			// This handler is called when rate limit is exceeded
			// We'll handle this in the Gin middleware instead
			w.WriteHeader(http.StatusTooManyRequests)
		}),
	)

	return func(c *gin.Context) {
		// Track if the request was rate limited
		rateLimited := false

		// Create a custom response writer to detect rate limiting
		originalWriter := c.Writer

		// Create a handler that will be called if rate limit is NOT exceeded
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If we reach here, rate limit was not exceeded
			c.Next()
		})

		// Apply the rate limiter
		limitedHandler := limiter(handler)

		// Capture the response to check if it was rate limited
		limitedHandler.ServeHTTP(c.Writer, c.Request)

		// Check if the response was a 429 (Too Many Requests)
		if c.Writer.Status() == http.StatusTooManyRequests {
			rateLimited = true
		}

		// If rate limited, provide a custom JSON response
		if rateLimited {
			c.Writer = originalWriter // Reset writer
			result := errors.HandleError(
				func() (interface{}, error) {
					return nil, errors.NewHTTPError(
						http.StatusTooManyRequests,
						"Rate limit exceeded because youre sending too many requests. Please try again later.",
						nil,
					)
				},
				"rate limit exceeded",
			)
			result.RespondWithJSON(c)
			c.Abort()
			return
		}

		// Add rate limit headers for client information
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Window", window.String())
	}
}