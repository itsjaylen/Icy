package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type clientInfo struct {
	count int
	reset time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*clientInfo)
)

func RateLimitMiddleware(next http.HandlerFunc, limit time.Duration, maxRequests int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Normalize IP (remove port)
		if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host
		}

		mu.Lock()
		info, exists := clients[ip]
		now := time.Now()

		if !exists || now.After(info.reset) {
			info = &clientInfo{
				count: 1,
				reset: now.Add(limit),
			}
			clients[ip] = info
		} else {
			if info.count >= maxRequests {
				mu.Unlock()
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			info.count++
		}
		mu.Unlock()

		next(w, r)
	}
}
