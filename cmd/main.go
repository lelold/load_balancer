package main

func main() {
	cfg, balancer, rateLimiter := initializeSystem()
	defer saveClients(rateLimiter)

	router := setupRouter(rateLimiter, balancer)
	server := startHTTPServer(cfg.Port, router)

	waitForShutdown(server)
}
