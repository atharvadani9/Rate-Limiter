package ratelimiter

import (
	"net/http"
	"strconv"
)

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}
		s := rl.Status(key)
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(s.Remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(s.ResetAt.Unix(), 10))

		if !rl.Allow(key) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
