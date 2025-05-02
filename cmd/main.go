// Package main запускает HTTP-сервер балансировщика нагрузки
package main

// main функция
func main() {
	cfg, balancer, rateLimiter := initializeSystem()
	defer saveClients(rateLimiter)

	router := setupRouter(rateLimiter, balancer)
	server := startHTTPServer(cfg.Port, router)

	waitForShutdown(server)
}
