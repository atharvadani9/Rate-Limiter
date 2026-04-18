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
	done     chan struct{}
}

func New(limit int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:    limit,
		duration: duration,
		store:    make(map[string]*window),
		done:     make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.done:
				return
			}
		}
	}()

	return rl
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

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for key, win := range rl.store {
		if win.expired(now, rl.duration) {
			delete(rl.store, key)
		}
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.done)
}
