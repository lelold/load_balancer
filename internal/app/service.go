// Package app реализует сервис, который объединяет функциональность лимитера и балансировщика
package app

import (
	"load_balancer/internal/domain/lb"
	"load_balancer/internal/domain/ratelimiter"
)

// Service инкапсулирует балансировщик и лимитер
type Service struct {
	Balancer lb.Balancer         // балансировщик
	Limiter  ratelimiter.Limiter // лимитер
}

// NewService создаёт новый сервис
func NewService(limiter ratelimiter.Limiter, balancer lb.Balancer) *Service {
	return &Service{
		Balancer: balancer,
		Limiter:  limiter,
	}
}
