// Package middleware provides middleware functions for analytics.
package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus metric collectors.
type Metrics struct {
	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	once         sync.Once
}

// NewMetrics initializes and registers Prometheus metrics.
func NewMetrics() *Metrics {
	metrics := &Metrics{}
	metrics.once.Do(func() {
		metrics.httpRequests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "status_code", "path"},
		)

		metrics.httpDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_duration_seconds",
				Help:    "Histogram of HTTP request durations",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "status_code", "path"},
		)

		prometheus.MustRegister(metrics.httpRequests)
		prometheus.MustRegister(metrics.httpDuration)
	})

	return metrics
}

// AnalyticsMiddleware captures request metrics for Prometheus.
func (m *Metrics) AnalyticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, req)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(ww.statusCode)

		m.httpRequests.WithLabelValues(req.Method, statusCode, req.URL.Path).Inc()
		m.httpDuration.WithLabelValues(req.Method, statusCode, req.URL.Path).Observe(duration)
	})
}

// MetricsHandler exposes the Prometheus metrics.
func (m *Metrics) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// responseWriterWrapper captures HTTP response status codes.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
