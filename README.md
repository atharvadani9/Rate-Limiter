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

## Development

```bash
go build ./...
go test ./...
go vet ./...
```
