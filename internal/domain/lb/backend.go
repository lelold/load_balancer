package lb

import (
	"load_balancer/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend описывает сервер для проксирования запросов
type Backend struct {
	URL               *url.URL               // адрес
	Alive             bool                   // флаг доступности
	ActiveConnections int                    // кол-во активных соединений
	Mutex             sync.RWMutex           // мьютекс
	ReverseProxy      *httputil.ReverseProxy // реверс прокси
}

// NewBackend создает Backend
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

	// обработчик ошибок reverse proxy помечает сервер как недоступный
	backend.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.SetAlive(false)
		logger.Infof("Backend %s failed: %v", backend.URL, err)
		http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
	}

	return backend
}

// IncConnections инкрементирует активные соединения к бэкенду
func (b *Backend) IncConnections() {
	b.Mutex.Lock()
	b.ActiveConnections++
	b.Mutex.Unlock()
}

// DecConnections декрементирует активные соединения к бэкенду
func (b *Backend) DecConnections() {
	b.Mutex.Lock()
	b.ActiveConnections--
	b.Mutex.Unlock()
}

// SetAlive устанавливает флаг Alive в "живой"
func (b *Backend) SetAlive(alive bool) {
	b.Mutex.Lock()
	b.Alive = alive
	b.Mutex.Unlock()
}

// IsAlive проверяет, жив ли сервер
func (b *Backend) IsAlive() bool {
	b.Mutex.RLock()
	defer b.Mutex.RUnlock()
	return b.Alive
}
