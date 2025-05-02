package handlers

import (
	"load_balancer/internal/app"
	"load_balancer/internal/domain/lb"
	"load_balancer/internal/domain/ratelimiter"
	"net/http"
)

// NewRouter создает новый роутер
func NewRouter(rateLimiter ratelimiter.Limiter, balancer lb.Balancer) http.Handler {
	mux := http.NewServeMux()

	clientService := app.NewService(rateLimiter, balancer)
	clientHandler := NewClientHandler(*clientService)

	mux.HandleFunc("/health", clientHandler.HealthHandler)
	mux.HandleFunc("/", clientHandler.ProxyHandler(clientHandler.svc))
	mux.HandleFunc("/clients", clientHandler.GetPostClients)
	mux.HandleFunc("/clients/", clientHandler.DeleteClient)

	return mux
}
