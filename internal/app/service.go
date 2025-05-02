package app

import (
	"load_balancer/internal/domain/lb"
	"load_balancer/internal/domain/ratelimiter"
)

type Service struct {
	Balancer lb.Balancer
	Limiter  ratelimiter.Limiter
}

func NewService(limiter ratelimiter.Limiter, balancer lb.Balancer) *Service {
	return &Service{
		Balancer: balancer,
		Limiter:  limiter,
	}
}
