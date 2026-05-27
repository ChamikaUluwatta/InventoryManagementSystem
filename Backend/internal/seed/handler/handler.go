package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/seed", h.Seed)
}

func (h *Handler) Seed(w http.ResponseWriter, r *http.Request) {
	_, _, err := h.service.Seed(r.Context())
	if err != nil {
		apperror.HandleError(w, apperror.Internal("seed failed", err))
		return
	}

	response := map[string]interface{}{
		"message": "Seed completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
