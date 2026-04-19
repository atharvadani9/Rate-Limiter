package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var dummyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestMiddleware_MissingAPIKey(t *testing.T) {
	rl := New(3, time.Minute)
	defer rl.Stop()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rl.Middleware(dummyHandler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMiddleware_AllowedRequest(t *testing.T) {
	rl := New(3, time.Minute)
	defer rl.Stop()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "test-key")
	rec := httptest.NewRecorder()

	rl.Middleware(dummyHandler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("X-RateLimit-Remaining") == "" {
		t.Fatal("expected X-RateLimit-Remaining header to be set")
	}
	if rec.Header().Get("X-RateLimit-Reset") == "" {
		t.Fatal("expected X-RateLimit-Reset header to be set")
	}
}

func TestMiddleware_RateLimitExceeded(t *testing.T) {
	rl := New(2, time.Minute)
	defer rl.Stop()

	handler := rl.Middleware(dummyHandler)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "test-key")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-API-Key", "test-key")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if rec.Header().Get("X-RateLimit-Remaining") == "" {
		t.Fatal("expected X-RateLimit-Remaining header on 429 response")
	}
	if rec.Header().Get("X-RateLimit-Reset") == "" {
		t.Fatal("expected X-RateLimit-Reset header on 429 response")
	}
}
