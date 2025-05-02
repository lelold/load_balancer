package ratelimiter

import (
	"sync"
	"time"
)

// Bucket реализует алгоритм token bucket, используется для контроля запросов от клиента
type Bucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewBucket конструирует новый бакет
func NewBucket(capacity, refillRate int) *Bucket {
	return &Bucket{
		capacity:   capacity,   // максимальное кол-во токенов
		tokens:     capacity,   // текущее кол-во токенов
		refillRate: refillRate, // токенов в секунду
		lastRefill: time.Now(), // время пополнения
	}
}

// Allow проверяет, можно ли пропустить запрос, и вычитает токен при разрешении
// возвращает true, если запрос разрешен
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
