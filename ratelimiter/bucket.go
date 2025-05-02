package ratelimiter

import (
	"sync"
	"time"
)

type Bucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewBucket(capacity, refillRate int) *Bucket {
	return &Bucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (b *Bucket) Allow() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	tokensToAdd := int(elapsed * float64(b.refillRate))

	if tokensToAdd > 0 {
		b.tokens = min(b.capacity, b.tokens+tokensToAdd)
		b.lastRefill = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}
