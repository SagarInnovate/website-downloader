package utils

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	mutex      sync.Mutex
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
// requestsPerSecond: number of requests allowed per second
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		tokens:     requestsPerSecond,
		maxTokens:  requestsPerSecond,
		refillRate: time.Second / time.Duration(requestsPerSecond),
		lastRefill: time.Now(),
	}
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	// Wait if no tokens available
	for rl.tokens <= 0 {
		rl.mutex.Unlock()
		time.Sleep(rl.refillRate)
		rl.mutex.Lock()

		now = time.Now()
		elapsed = now.Sub(rl.lastRefill)
		tokensToAdd = int(elapsed / rl.refillRate)

		if tokensToAdd > 0 {
			rl.tokens += tokensToAdd
			if rl.tokens > rl.maxTokens {
				rl.tokens = rl.maxTokens
			}
			rl.lastRefill = now
		}
	}

	// Consume a token
	rl.tokens--
}

// TryAcquire attempts to acquire a token without blocking
// Returns true if a token was acquired, false otherwise
func (rl *RateLimiter) TryAcquire() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	// Try to consume a token
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}
