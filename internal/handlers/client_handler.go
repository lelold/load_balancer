package handlers

import (
	"encoding/json"
	"load_balancer/internal/app"
	"load_balancer/internal/logger"
	"net/http"
	"strings"
)

type ClientConfig struct {
	ClientID   string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RatePerSec int    `json:"rate_per_sec"`
}

type ClientHandler struct {
	svc app.Service
}

func NewClientHandler(service app.Service) *ClientHandler {
	return &ClientHandler{svc: service}
}

func (h *ClientHandler) CreateOrUpdateClient(w http.ResponseWriter, r *http.Request) {
	var cfg ClientConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	logger.Infof("Created or updated client: %s", cfg.ClientID)
	h.svc.Limiter.AddClient(cfg.ClientID, cfg.Capacity, cfg.RatePerSec)
	w.WriteHeader(http.StatusOK)
}

func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	clients := h.svc.Limiter.ListClients()
	resp, _ := json.Marshal(clients)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *ClientHandler) GetPostClients(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateOrUpdateClient(w, r)
	case http.MethodGet:
		h.GetClients(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		id := strings.TrimPrefix(r.URL.Path, "/clients/")
		if id == "" {
			http.Error(w, `{"code":400,"message":"Missing client ID"}`, http.StatusBadRequest)
			return
		}
		if ok := h.svc.Limiter.DeleteClient(id); ok {
			logger.Infof("Deleted client: %s", id)
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, `{"code":404,"message":"Client not found"}`, http.StatusNotFound)
		}
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *ClientHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *ClientHandler) ProxyHandler(svc app.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := r.Header.Get("X-Client-ID")
		if clientID == "" {
			http.Error(w, "Missing X-Client-ID", http.StatusBadRequest)
			return
		}

		if !svc.Limiter.Allow(clientID) {
			logger.Infof("Rate limit exceeded for %s", clientID)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		svc.Balancer.ServeHTTP(w, r)
	}
}
