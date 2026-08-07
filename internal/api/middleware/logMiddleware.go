package middleware

import (
	"plantao/internal/domain/log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(log log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqLog := log.With(
			"method", c.Request.Method,
			"path", c.FullPath(),
			"ip", c.ClientIP(),
			"request_id", c.GetHeader("X-Request-ID"),
		)

		reqLog.Info("request started")

		c.Next() // executa o handler

		reqLog.Info("request completed",
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)

		// loga erros que ocorreram durante o handler
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				reqLog.Error("handler error", "error", e.Error())
			}
		}
	}
}
