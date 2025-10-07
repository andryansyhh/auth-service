package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLogger() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)

		logger.Info("request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"latency", latency.String(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
