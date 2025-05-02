package lb

import "net/http"

// Интерфейс балансировщика
type Balancer interface {
	http.Handler
}

// Интерфейс стратегии
type Strategy interface {
	GetBackend([]*Backend) *Backend
}
