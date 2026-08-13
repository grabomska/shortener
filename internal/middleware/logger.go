package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.Infow("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"duration", time.Since(start).String(),
			"status", c.Writer.Status(),
			"bytes", c.Writer.Size())
	}
}
