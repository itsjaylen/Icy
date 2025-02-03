package middleware

import (
	"net/http"
	"sync"
	"time"
)

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
