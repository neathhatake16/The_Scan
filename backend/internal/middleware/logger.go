package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/pkg/logger"
)

// RequestLogger logs every HTTP request with method, path, status, latency.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.Infow("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
		)
	}
}