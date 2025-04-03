package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiterStore manages request rate limits per client.
type RateLimiterStore struct {
	store sync.Map
}

var globalRateLimiter = &RateLimiterStore{}

// rateLimiter keeps track of request counts and reset time.
type rateLimiter struct {
	visits map[string]int
	reset  time.Time
	sync.Mutex
}

// RateLimiter limits requests within a given duration.
func RateLimiter(next http.Handler, limit int, duration time.Duration) http.Handler {
	rl := &rateLimiter{
		visits: make(map[string]int),
		reset:  time.Now().Add(duration),
	}

	go func() {
		for {
			time.Sleep(duration)
			rl.Lock()
			rl.visits = make(map[string]int)
			rl.reset = time.Now().Add(duration)
			rl.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.Lock()
		defer rl.Unlock()

		ip := r.RemoteAddr
		if rl.visits[ip] >= limit {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}

		rl.visits[ip]++
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware applies a request rate limit based on time duration.
func RateLimitMiddleware(next http.HandlerFunc, limit time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if lastRequest, exists := globalRateLimiter.store.Load(clientIP); exists {
			if lastReqTime, ok := lastRequest.(time.Time); ok && time.Since(lastReqTime) < limit {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}
		}
		globalRateLimiter.store.Store(clientIP, time.Now())
		next(w, r)
	}
}
