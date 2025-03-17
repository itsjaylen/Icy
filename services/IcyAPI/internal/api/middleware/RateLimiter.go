package middleware

import (
	"net/http"
	"sync"
	"time"
)

var rateLimitStore     sync.Map

// rateLimiter tracks requests per client
type rateLimiter struct {
	sync.Mutex
	visits map[string]int
	reset  time.Time
}

// RateLimiter middleware limits requests within a given duration
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
			return
		}

		rl.visits[ip]++
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.HandlerFunc, limit time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if lastRequest, exists := rateLimitStore.Load(clientIP); exists && time.Since(lastRequest.(time.Time)) < limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		rateLimitStore.Store(clientIP, time.Now())
		next(w, r)
	}
}