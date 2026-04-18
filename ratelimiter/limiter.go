package ratelimiter

import (
	"sync"
	"time"
)

type Status struct {
	Remaining int
	ResetAt   time.Time
}

type RateLimiter struct {
	limit    int
	duration time.Duration
	mu       sync.RWMutex
	store    map[string]*window
}

func New(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		duration: duration,
		store:    make(map[string]*window),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	win, ok := rl.store[key]
	if !ok {
		rl.store[key] = &window{counter: 1, startedAt: now}
		return true
	}
	if win.expired(now, rl.duration) {
		win.reset(now)
		return true
	}
	if win.counter < rl.limit {
		win.add()
		return true
	}
	return false
}

func (rl *RateLimiter) Status(key string) Status {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	win, ok := rl.store[key]
	if !ok {
		return Status{
			Remaining: rl.limit,
			ResetAt:   time.Now().Add(rl.duration),
		}
	}

	if win.expired(time.Now(), rl.duration) {
		return Status{
			Remaining: rl.limit,
			ResetAt:   time.Now().Add(rl.duration),
		}
	} else {
		return Status{
			Remaining: rl.limit - win.counter,
			ResetAt:   time.Now().Add(rl.duration),
		}
	}
}
