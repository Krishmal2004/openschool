package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger replaces gin.Default()'s plain-text access log, which
// prints the full request line — including query strings, so search terms,
// names and index numbers end up in plain-text logs (S10). It logs the
// route template rather than the raw path, so /students/:id doesn't fan out
// into one log shape per student.
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		slog.Info("request",
			"method", c.Request.Method,
			"route", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"request_id", RequestIDFromContext(c),
			"user_id", c.GetString("userID"),
		)
	}
}
