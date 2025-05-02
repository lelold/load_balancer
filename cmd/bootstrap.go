package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"load_balancer/internal/config"
	"load_balancer/internal/domain/lb"
	"load_balancer/internal/domain/ratelimiter"
	"load_balancer/internal/handlers"
	"load_balancer/internal/logger"
)

// initializeSystem загружает конфигурацию, инициализирует стратегию балансировки, балансировщик и лимитер
func initializeSystem() (*config.Config, lb.Balancer, ratelimiter.Limiter) {
	cfgPath := "assets/config.json"
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

	balancer := lb.NewLoadBalancer(cfg.Backends, strategy)

	rateLimiter := ratelimiter.NewRateLimiter(cfg.RateLimit.Capacity, cfg.RateLimit.RefillRate)

	if err := rateLimiter.LoadClientsFromFile("assets/clients.json"); err != nil {
		logger.Infof("Could not load clients: %v", err)
	}

	return cfg, balancer, rateLimiter
}

// saveClients сохраняетов клиентов в файл assets/clients.json
func saveClients(limiter ratelimiter.Limiter) {
	if err := limiter.SaveClientsToFile("assets/clients.json"); err != nil {
		logger.Errorf("Could not save clients: %v", err)
	}
}

// setupRouter инициализирует роутер
func setupRouter(limiter ratelimiter.Limiter, balancer lb.Balancer) http.Handler {
	return handlers.NewRouter(limiter, balancer)
}

// startHTTPServer запускает сервер на заданном порту
func startHTTPServer(port string, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		logger.Infof("Listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
		}
	}()

	return server
}

// waitForShutdown асинхронно ожидает завершения
func waitForShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

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
