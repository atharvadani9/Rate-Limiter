# Rate Limiter

A fixed-window, in-memory rate limiter library written in Go, keyed by API key.

## Usage

```go
import ratelimiter "github.com/atharvadani9/Rate-Limiter"

rl := ratelimiter.New(100, time.Minute) // 100 requests per minute

if rl.Allow("api-key-123") {
    // handle request
} else {
    // reject request
}

status := rl.Status("api-key-123")
fmt.Println(status.Remaining) // requests left in current window
fmt.Println(status.ResetAt)   // when the window resets
```

## API

### `New(limit int, windowDuration time.Duration) *RateLimiter`
Creates a new rate limiter.

### `Allow(key string) bool`
Returns `true` if the request is within the rate limit and increments the counter.

### `Status(key string) Status`
Returns the current rate limit state without consuming quota.

```go
type Status struct {
    Remaining int
    ResetAt   time.Time
}
```

## HTTP Middleware

```go
rl := ratelimiter.New(3, 5*time.Second)
defer rl.Stop()

mux := http.NewServeMux()
mux.HandleFunc("/", handler)
http.ListenAndServe(":8080", rl.Middleware(mux))
```

Requests must include an `X-API-Key` header. The middleware sets the following headers on every response:
- `X-RateLimit-Remaining` — requests left in the current window
- `X-RateLimit-Reset` — unix timestamp of when the window resets

**Testing with curl:**
```bash
# valid request
curl -i -H "X-API-Key: mykey" http://localhost:8080/

# trigger rate limit (runs 4 times)
for i in {1..4}; do curl -i -H "X-API-Key: mykey" http://localhost:8080/; done

# missing key - returns 401
curl -i http://localhost:8080/
```

## Development

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
```
