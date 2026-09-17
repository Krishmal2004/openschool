package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// requestsTotal and requestDuration are the two series the observability
// checklist asks for: error rate (via the status label) and p95 latency
// (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 5). Registered once at
// package init so repeated calls to Metrics() in tests don't panic on a
// duplicate registration.
var (
	requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "openschool_http_requests_total",
		Help: "Total HTTP requests, labelled by method, route template and status.",
	}, []string{"method", "route", "status"})

	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "openschool_http_request_duration_seconds",
		Help:    "HTTP request latency in seconds, labelled by method and route template.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration)
}

// Metrics records each request's count and latency under its route
// template (never the raw path, so /students/:id doesn't fan out into one
// series per student).
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		requestsTotal.WithLabelValues(c.Request.Method, route, strconv.Itoa(c.Writer.Status())).Inc()
		requestDuration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}
