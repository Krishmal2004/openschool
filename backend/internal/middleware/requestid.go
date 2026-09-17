package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openschool-org/openschool/internal/apierror"
)

// RequestIDHeader is the header a per-request id travels on, both back to
// the client and, if the proxy already assigned one upstream, in from it.
const RequestIDHeader = "X-Request-Id"

// RequestID assigns a per-request id — reusing one already set by the proxy
// when present — so a single request can be traced across structured logs
// and generic error responses (S4, S10).
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
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
