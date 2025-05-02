package lb

import (
	"load_balancer/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL               *url.URL
	Alive             bool
	ActiveConnections int
	Mutex             sync.RWMutex
	ReverseProxy      *httputil.ReverseProxy
}

func NewBackend(rawURL string) *Backend {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Errorf("Invalid backend URL: %s", rawURL)
		return nil
	}

	backend := &Backend{
		URL:          parsedURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}

	backend.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.SetAlive(false)
		logger.Infof("Backend %s failed: %v", backend.URL, err)
		http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
	}

	return backend
}

func (b *Backend) IncConnections() {
	b.Mutex.Lock()
	b.ActiveConnections++
	b.Mutex.Unlock()
}

func (b *Backend) DecConnections() {
	b.Mutex.Lock()
	b.ActiveConnections--
	b.Mutex.Unlock()
}

func (b *Backend) SetAlive(alive bool) {
	b.Mutex.Lock()
	b.Alive = alive
	b.Mutex.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()
	return b.Alive
}
