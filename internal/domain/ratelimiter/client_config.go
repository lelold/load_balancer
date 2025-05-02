package ratelimiter

import (
	"encoding/json"
	"os"
	"time"
)

// ClientConfig описывает конфигурацию клиента
type ClientConfig struct {
	ClientID   string `json:"client_id"`    // уникальный идентификатор клиента
	Capacity   int    `json:"capacity"`     // максимальное количество токенов
	RatePerSec int    `json:"rate_per_sec"` // скорость пополнения токенов в секунду
	Tokens     int    `json:"tokens"`       // текущее количество токенов
	LastRefill int64  `json:"last_refill"`  // время последнего пополнения
}

// LoadClientsFromFile загружает конфигурации клиентов из json-файла
func (r *RateLimiter) LoadClientsFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var clients []ClientConfig
	if err := json.NewDecoder(file).Decode(&clients); err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, c := range clients {
		b := &Bucket{
			capacity:   c.Capacity,
			tokens:     c.Tokens,
			refillRate: c.RatePerSec,
			lastRefill: time.Unix(c.LastRefill, 0),
		}
		r.buckets[c.ClientID] = b
	}
	return nil
}

// SaveClientsToFile сохраняет конфигурации клиентов из json-файла
func (r *RateLimiter) SaveClientsToFile(path string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var clients []ClientConfig
	for id, b := range r.buckets {
		b.mutex.Lock()
		clients = append(clients, ClientConfig{
			ClientID:   id,
			Capacity:   b.capacity,
			RatePerSec: b.refillRate,
			Tokens:     b.tokens,
			LastRefill: b.lastRefill.Unix(),
		})
		b.mutex.Unlock()
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(clients)
}
