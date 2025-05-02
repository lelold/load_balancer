// Package lb реализует балансировщик с тремя стратегиями (round-robin, least connections, random)
package lb

import (
	"load_balancer/internal/logger"
	"net/http"
	"time"
)

// LoadBalancer описывает балансировщик с заданной стратегией и списком серверов
type LoadBalancer struct {
	backends []*Backend
	strategy Strategy
}

// NewLoadBalancer инициализирует балансировщик и проверяет, живы ли сервера
func NewLoadBalancer(backendUrls []string, strategy Strategy) *LoadBalancer {
	backends := make([]*Backend, 0, len(backendUrls))
	for _, addr := range backendUrls {
		backend := NewBackend(addr)
		if backend != nil {
			backends = append(backends, backend)
		}
	}

	lb := &LoadBalancer{
		backends: backends,
		strategy: strategy,
	}

	go lb.healthCheck()
	return lb
}

// ServeHTTP обрабатывает http запрос и выбирает сервер по стратегии, учитывая число активных соединений
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.strategy.GetBackend(lb.backends)
	if backend == nil {
		http.Error(w, "No available backend", http.StatusServiceUnavailable)
		return
	}

	backend.IncConnections()
	defer backend.DecConnections()

	logger.Infof("Forwarding request to %s", backend.URL)
	backend.ReverseProxy.ServeHTTP(w, r)
}

// healthCheck периодически проверяет доступность каждого сервера, отправляя запрос на endpoint /health
func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(10 * time.Second) // каждые 10 секунд
	for range ticker.C {
		for _, b := range lb.backends {
			go func(backend *Backend) {
				resp, err := http.Get(backend.URL.String() + "/health")
				if err != nil || resp.StatusCode != http.StatusOK { // если /health не отвечает, то говорим, что сервер умер
					backend.SetAlive(false)
					logger.Infof("Health check failed for %s", backend.URL)
				} else {
					logger.Infof("Health check passed for %s", backend.URL)
					backend.SetAlive(true)
				}
			}(b)
		}
	}
}
