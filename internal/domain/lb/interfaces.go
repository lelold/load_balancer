package lb

import "net/http"

type Balancer interface {
	http.Handler
}

type Strategy interface {
	GetBackend([]*Backend) *Backend
}
