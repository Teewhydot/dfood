package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"dfood/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("Request processed",
			"client_id", param.ClientIP,
			"timestamp", param.TimeStamp.Format(time.RFC1123),
			"method", param.Method,
			"path", param.Path,
			"protocol", param.Request.Proto,
			"status_code", param.StatusCode,
			"latency", param.Latency,
			"user_agent", param.Request.UserAgent(),
		)
		return ""
	})
}