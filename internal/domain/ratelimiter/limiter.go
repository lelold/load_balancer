package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	buckets       map[string]*Bucket
	mutex         sync.RWMutex
	defaultCap    int
	defaultRefill int
}

func NewRateLimiter(capacity, refill int) *RateLimiter {
	return &RateLimiter{
		buckets:       make(map[string]*Bucket),
		defaultCap:    capacity,
		defaultRefill: refill,
	}
}

func (r *RateLimiter) getBucket(clientID string) *Bucket {
	r.mutex.RLock()
	bucket, exists := r.buckets[clientID]
	r.mutex.RUnlock()

	if exists {
		return bucket
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if bucket, exists := r.buckets[clientID]; exists {
		return bucket
	}

	newBucket := NewBucket(r.defaultCap, r.defaultRefill)
	r.buckets[clientID] = newBucket
	return newBucket
}

func (r *RateLimiter) Allow(clientID string) bool {
	return r.getBucket(clientID).Allow()
}

func (r *RateLimiter) AddClient(clientID string, capacity, refill int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buckets[clientID] = NewBucket(capacity, refill)
}

func (r *RateLimiter) UpdateClient(clientID string, capacity, refill int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if bucket, exists := r.buckets[clientID]; exists {
		bucket.capacity = capacity
		bucket.refillRate = refill
		bucket.tokens = capacity
		bucket.lastRefill = time.Now()
	}
}

func (r *RateLimiter) ListClients() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	clients := make([]string, 0, len(r.buckets))
	for id := range r.buckets {
		clients = append(clients, id)
	}
	return clients
}

func (r *RateLimiter) DeleteClient(id string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.buckets[id]; exists {
		delete(r.buckets, id)
		return true
	}
	return false
}
