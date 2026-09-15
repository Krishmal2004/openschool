// Package httpx contains transport helpers shared by HTTP modules.
package httpx

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// BindStrict decodes and validates JSON while rejecting unknown fields.
func BindStrict(c *gin.Context, destination any) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(destination)
}
