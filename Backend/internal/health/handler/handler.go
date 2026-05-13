package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/health/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.healthCheck)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	results, status := h.service.Check(r.Context())

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"status": status,
		"checks": results,
	})
}
