# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...          # build
go test ./...           # run all tests
go test -run TestName   # run a single test
go vet ./...            # lint
```

## Architecture

A fixed-window, in-memory rate limiter library written in Go, keyed by API key.

**Module:** `github.com/atharvadani9/rate-limiter`

**Core types:**
- `window` (window.go) — tracks request count and window start time per key
- `RateLimiter` (limiter.go) — holds config (limit, duration) and the in-memory store; public API is `Allow(key string) bool` and `Status(key string) Status`

**Key behaviors:**
- `Allow()` — increments counter if within limit; resets window if expired; denies otherwise
- `Status()` — read-only, never increments counter; returns full quota for unseen keys
- Store is protected by `sync.RWMutex` — reads via `RLock`, writes via `Lock`
- A background goroutine in `New()` periodically evicts expired windows to prevent memory leaks
