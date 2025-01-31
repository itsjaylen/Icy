package middleware

// AnalyticsMiddleware captures request metrics for Prometheus
// Experimental

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

var (
	// Define a Prometheus counter for request count
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gin_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status_code", "path"},
	)

	// Define a Prometheus histogram for request duration
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gin_http_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets, // You can customize the buckets
		},
		[]string{"method", "status_code", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

// AnalyticsMiddleware captures request metrics for Prometheus
func AnalyticsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())

		// Record the metrics
		httpRequests.WithLabelValues(c.Request.Method, statusCode, c.Request.URL.Path).Inc()
		httpDuration.WithLabelValues(c.Request.Method, statusCode, c.Request.URL.Path).
			Observe(duration)
	}
}

// MetricsHandler exposes the Prometheus metrics
func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
