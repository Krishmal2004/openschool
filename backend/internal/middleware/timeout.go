package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestTimeout bounds every request's context to d, so one slow downstream
// call (database, mailer, ThunderID) can't hold a handler goroutine and its
// pooled database connection indefinitely (S15). Matches the http.Server's
// own read/write timeouts in cmd/api/main.go.
func RequestTimeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
