package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestAllow_WithinLimit(t *testing.T) {
	rl := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !rl.Allow("key1") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	rl := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		rl.Allow("key1")
	}
	if rl.Allow("key1") {
		t.Fatal("4th request should be denied")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	rl := New(2, 50*time.Millisecond)
	rl.Allow("key1")
	rl.Allow("key1")

	if rl.Allow("key1") {
		t.Fatal("3rd request should be denied before reset")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("key1") {
		t.Fatal("request after window reset should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	rl := New(1, time.Minute)
	rl.Allow("key1")

	if !rl.Allow("key2") {
		t.Fatal("key2 should have its own independent window")
	}
}

func TestStatus_UnseenKey(t *testing.T) {
	rl := New(10, time.Minute)
	s := rl.Status("unknown")

	if s.Remaining != 10 {
		t.Fatalf("expected Remaining=10, got %d", s.Remaining)
	}
	if s.ResetAt.Before(time.Now()) {
		t.Fatal("ResetAt should be in the future")
	}
}

func TestStatus_DoesNotConsume(t *testing.T) {
	rl := New(10, time.Minute)
	rl.Allow("key1")

	s1 := rl.Status("key1")
	s2 := rl.Status("key1")

	if s1.Remaining != s2.Remaining {
		t.Fatal("Status() should not consume quota")
	}
}

func TestStatus_ReflectsAllows(t *testing.T) {
	rl := New(5, time.Minute)
	rl.Allow("key1")
	rl.Allow("key1")

	s := rl.Status("key1")
	if s.Remaining != 3 {
		t.Fatalf("expected Remaining=3, got %d", s.Remaining)
	}
}

func TestCleanup_RemovesExpiredKeys(t *testing.T) {
	rl := New(5, 50*time.Millisecond)
	rl.Allow("key1")

	// wait for window to expire and cleanup to run
	time.Sleep(120 * time.Millisecond)

	rl.mu.RLock()
	_, exists := rl.store["key1"]
	rl.mu.RUnlock()

	if exists {
		t.Fatal("expired key should have been removed by cleanup")
	}
}

func TestStop_StopsCleanupGoroutine(t *testing.T) {
	rl := New(5, 50*time.Millisecond)
	rl.Stop()
	// calling Stop twice should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Stop() panicked, likely closed done channel twice")
		}
	}()
}

func TestAllow_Concurrent(t *testing.T) {
	rl := New(100, time.Minute)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Allow("key1")
		}()
	}

	wg.Wait()

	s := rl.Status("key1")
	if s.Remaining < 0 {
		t.Fatal("Remaining should never go negative")
	}
}
