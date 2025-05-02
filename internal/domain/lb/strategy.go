package lb

import (
	"math/rand"
	"sync"
)

// RoundRobinStrategy реализует round-robin стратегию
type RoundRobinStrategy struct {
	current int
	mutex   sync.Mutex
}

func (rr *RoundRobinStrategy) GetBackend(backends []*Backend) *Backend {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	count := len(backends)
	for i := 0; i < count; i++ {
		idx := (rr.current + i) % count
		b := backends[idx]
		if b.IsAlive() { // учитываем, жив ли сервер
			rr.current = (idx + 1) % count
			return b
		}
	}
	return nil
}

// RandomStrategy реализует random стратегию
type RandomStrategy struct{}

func (r *RandomStrategy) GetBackend(backends []*Backend) *Backend {
	live := []*Backend{}
	for _, b := range backends {
		if b.IsAlive() {
			live = append(live, b)
		}
	}
	if len(live) == 0 {
		return nil
	}
	return live[rand.Intn(len(live))]
}

// LeastConnectionsStrategy реализует least connections стратегию
type LeastConnectionsStrategy struct{}

func (l *LeastConnectionsStrategy) GetBackend(backends []*Backend) *Backend {
	var best *Backend
	min := int(^uint(0) >> 1)
	for _, b := range backends {
		if b.IsAlive() {
			b.Mutex.RLock()
			conns := b.ActiveConnections
			b.Mutex.RUnlock()
			if conns < min {
				best = b
				min = conns
			}
		}
	}
	return best
}
