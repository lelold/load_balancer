package ratelimiter

// Интерфейс лимитера
type Limiter interface {
	Allow(clientID string) bool
	AddClient(clientID string, capacity, refill int)
	UpdateClient(clientID string, capacity, refill int)
	ListClients() []string
	DeleteClient(clientID string) bool
	LoadClientsFromFile(path string) error
	SaveClientsToFile(path string) error
}
