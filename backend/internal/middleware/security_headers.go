package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets a small set of defensive headers appropriate for a
// JSON-only API (no HTML is ever served here, so no CSP is needed — the SPA
// sets its own at the proxy, see docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md S7).
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "same-origin")
		// Belt-and-braces alongside the proxy's own HSTS header (S7): harmless
		// over plain HTTP in local development, since browsers ignore
		// Strict-Transport-Security on a non-HTTPS response.
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Next()
	}
}
