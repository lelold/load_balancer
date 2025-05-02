package testutils

import (
	"fmt"
	"net/http"

	"load_balancer/internal/app"
	"load_balancer/internal/domain/lb"
	"load_balancer/internal/domain/ratelimiter"
	"load_balancer/internal/handlers"
)

func NewHandlerWithMockedDeps() http.Handler {
	backends := []string{"http://localhost:8081", "http://localhost:8082"}

	strategy := &lb.RoundRobinStrategy{}

	var balancer lb.Balancer = lb.NewLoadBalancer(backends, strategy)

	var rateLimiter ratelimiter.Limiter = ratelimiter.NewRateLimiter(5, 5)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("X-Client-ID")
		if clientID == "" {
			http.Error(w, "Missing X-Client-ID", http.StatusBadRequest)
			return
		}

		if !rateLimiter.Allow(clientID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		balancer.ServeHTTP(w, r)
	})

	clientService := app.NewService(rateLimiter, balancer)
	clientHandler := handlers.NewClientHandler(*clientService)

	mux.HandleFunc("/clients", clientHandler.GetPostClients)
	mux.HandleFunc("/clients/", clientHandler.DeleteClient)

	return mux
}
