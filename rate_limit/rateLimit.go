package ratelimit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int
	tokens   float64
	rate     float64 // tokens per second
	last     time.Time
	mu       sync.Mutex
}

func NewTokenBucket(rate float64, capacity int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   float64(capacity),
		rate:     rate,
		last:     time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()

	// refill tokens
	tb.tokens += elapsed * tb.rate
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	tb.last = now

	// check if request can pass
	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}

	return false
}
