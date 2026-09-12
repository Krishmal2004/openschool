package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// bindStrict decodes and validates a JSON request body, rejecting any field not defined on the target struct.
func bindStrict(c *gin.Context, obj any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(obj); err != nil {
		return err
	}
	return binding.Validator.ValidateStruct(obj)
}
