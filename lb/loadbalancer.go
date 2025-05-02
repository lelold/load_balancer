package lb

import (
	"load_balancer/internal/logger"
	"net/http"
	"time"
)

type Balancer interface {
	http.Handler
}

type LoadBalancer struct {
	backends []*Backend
	strategy Strategy
}

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

func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		for _, b := range lb.backends {
			go func(backend *Backend) {
				resp, err := http.Get(backend.URL.String() + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
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
