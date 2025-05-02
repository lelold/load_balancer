// Package ratelimiter реализует лимитер с алгоритмом token bucket для ограничения запросов на клиента
package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter управляет наборами токен-бакетов для разных клиентов
type RateLimiter struct {
	buckets       map[string]*Bucket // мапа клиентских id к их бакетам
	mutex         sync.RWMutex       // защищает доступ к бакетам
	defaultCap    int                // ёмкость по умолчанию
	defaultRefill int                // скорость пополнения по умолчанию
}

// NewRateLimiter создает новый лимитер
func NewRateLimiter(capacity, refill int) *RateLimiter {
	return &RateLimiter{
		buckets:       make(map[string]*Bucket),
		defaultCap:    capacity,
		defaultRefill: refill,
	}
}

// getBucket возвращает новый бакет клиента или создает новый
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

// Allow проверяет, может ли клиент выполнить запрос
func (r *RateLimiter) Allow(clientID string) bool {
	return r.getBucket(clientID).Allow()
}

// AddClient добавляет нового клиента по id
func (r *RateLimiter) AddClient(clientID string, capacity, refill int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buckets[clientID] = NewBucket(capacity, refill)
}

// UpdateClient обновляет конфигурацию клиента по id
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

// ListClients возвращает список всех клиентов
func (r *RateLimiter) ListClients() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	clients := make([]string, 0, len(r.buckets))
	for id := range r.buckets {
		clients = append(clients, id)
	}
	return clients
}

// DeleteClient удаляет клиента по id
func (r *RateLimiter) DeleteClient(id string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.buckets[id]; exists {
		delete(r.buckets, id)
		return true
	}
	return false
}
