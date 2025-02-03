package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Define a Prometheus counter for request count
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status_code", "path"},
	)

	// Define a Prometheus histogram for request duration
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets, 
		},
		[]string{"method", "status_code", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(httpDuration)
}

// AnalyticsMiddleware captures request metrics for Prometheus
func AnalyticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		// Calculate duration
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(ww.statusCode)

		// Record the metrics
		httpRequests.WithLabelValues(r.Method, statusCode, r.URL.Path).Inc()
		httpDuration.WithLabelValues(r.Method, statusCode, r.URL.Path).Observe(duration)
	})
}

// MetricsHandler exposes the Prometheus metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// responseWriterWrapper is used to capture HTTP response status codes
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
