package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsRecordsRouteTemplateNotRawPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Metrics())
	router.GET("/students/:id", func(c *gin.Context) { c.Status(200) })

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/students/aaa", nil))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/students/bbb", nil))

	got := testutil.ToFloat64(requestsTotal.WithLabelValues("GET", "/students/:id", "200"))
	if got != 2 {
		t.Fatalf("requestsTotal for /students/:id = %v, want 2 (raw student ids must not create separate series)", got)
	}
}

func TestMetricsLabelsUnmatchedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Metrics())
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/does-not-exist", nil))

	got := testutil.ToFloat64(requestsTotal.WithLabelValues("GET", "unmatched", "404"))
	if got < 1 {
		t.Fatalf("requestsTotal for unmatched route = %v, want >= 1", got)
	}
}
