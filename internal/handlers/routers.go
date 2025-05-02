package handlers

import (
	"load_balancer/lb"
	"load_balancer/ratelimiter"
	"net/http"
)

func NewRouter(rateLimiter ratelimiter.Limiter, balancer lb.Balancer) http.Handler {
	mux := http.NewServeMux()

	clientHandler := &ClientHandler{Limiter: rateLimiter}

	mux.HandleFunc("/health", clientHandler.HealthHandler)
	mux.HandleFunc("/", clientHandler.ProxyHandler(rateLimiter, balancer))
	mux.HandleFunc("/clients", clientHandler.GetPostClients)
	mux.HandleFunc("/clients/", clientHandler.DeleteClient)

	return mux
}
