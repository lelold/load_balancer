package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"load_balancer/config"
	"load_balancer/internal/handlers"
	"load_balancer/internal/logger"
	"load_balancer/lb"
	"load_balancer/ratelimiter"
)

func main() {
	cfgPath := "config.json"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	var strategy lb.Strategy
	switch cfg.Strategy {
	case "least_connections":
		strategy = &lb.LeastConnectionsStrategy{}
	case "random":
		strategy = &lb.RandomStrategy{}
	default:
		strategy = &lb.RoundRobinStrategy{}
	}

	logger.Info("Starting Load Balancer...")

	var balancer lb.Balancer = lb.NewLoadBalancer(cfg.Backends, strategy)

	const clientsFile = "clients.json"

	var rateLimiter ratelimiter.Limiter = ratelimiter.NewRateLimiter(cfg.RateLimit.Capacity, cfg.RateLimit.RefillRate)

	if err := rateLimiter.LoadClientsFromFile(clientsFile); err != nil {
		logger.Infof("Could not load clients: %v", err)
	}

	defer func() {
		if err := rateLimiter.SaveClientsToFile(clientsFile); err != nil {
			logger.Errorf("Could not save clients: %v", err)
		}
	}()

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
			logger.Infof("Rate limit exceeded for %s", clientID)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		balancer.ServeHTTP(w, r)
	})

	clientHandler := &handlers.ClientHandler{Limiter: rateLimiter}

	mux.HandleFunc("/clients", clientHandler.GetPostClients)

	mux.HandleFunc("/clients/", clientHandler.DeleteClient)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
		}
	}()

	<-stop
	logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Graceful shutdown failed: %v", err)
	} else {
		logger.Info("Server stopped cleanly.")
	}
}
