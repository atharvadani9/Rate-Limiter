package ratelimiter

import "time"

type window struct {
	counter   int
	startedAt time.Time
}

func (w *window) expired(now time.Time, windowDuration time.Duration) bool {
	return !now.Before(w.startedAt.Add(windowDuration))
}

func (w *window) reset(now time.Time) {
	w.counter = 1
	w.startedAt = now
}

func (w *window) add() {
	w.counter++
}
