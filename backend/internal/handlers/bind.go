package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/openschool-org/openschool/internal/platform/httpx"
)

// bindStrict decodes and validates a JSON request body, rejecting any field not defined on the target struct.
func bindStrict(c *gin.Context, obj any) error {
	return httpx.BindStrict(c, obj)
}
