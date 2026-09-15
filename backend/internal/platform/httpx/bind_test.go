package httpx

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBindStrictRejectsUnknownFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Grade 1","unknown":true}`))
	request.Header.Set("Content-Type", "application/json")
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request

	var destination struct {
		Name string `json:"name" binding:"required"`
	}
	if err := BindStrict(context, &destination); err == nil {
		t.Fatal("BindStrict() accepted an unknown field")
	}
}
