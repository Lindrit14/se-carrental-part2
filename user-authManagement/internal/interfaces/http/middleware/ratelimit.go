package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimit is a tiny in-memory IP-based limiter.
// For multi-replica deployments switch to a shared store (Redis); this is
// fine for single-instance dev/local and for ACA Min=Max=1.
type RateLimit struct {
	max    int           // requests per window
	window time.Duration // window length

	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count       int
	windowStart time.Time
}

func NewRateLimit(maxPerMin int) *RateLimit {
	return &RateLimit{
		max:     maxPerMin,
		window:  time.Minute,
		buckets: make(map[string]*bucket),
	}
}

func (rl *RateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ipKey(r)
		if !rl.allow(ip, time.Now()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code":    "rate_limited",
				"message": "too many requests",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimit) allow(key string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[key]
	if !ok || now.Sub(b.windowStart) > rl.window {
		rl.buckets[key] = &bucket{count: 1, windowStart: now}
		return true
	}
	if b.count >= rl.max {
		return false
	}
	b.count++
	return true
}

func ipKey(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return v
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
