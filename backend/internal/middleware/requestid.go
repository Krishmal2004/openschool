package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
)

// RequestIDHeader is the header a per-request id is returned on.
const RequestIDHeader = "X-Request-Id"

// RequestID assigns a per-request id so a single request can be traced
// across structured logs and generic error responses (S4, S10). The id is
// always generated here rather than trusted from an inbound header: a
// caller-supplied X-Request-Id would otherwise be echoed straight back and
// recorded in logs, making correlation ambiguous.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set(apierror.RequestIDContextKey, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}

// RequestIDFromContext returns the current request's id, or "" if RequestID didn't run.
func RequestIDFromContext(c *gin.Context) string {
	id, _ := c.Get(apierror.RequestIDContextKey)
	s, _ := id.(string)
	return s
}
