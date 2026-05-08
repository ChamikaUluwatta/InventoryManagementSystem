package seed

import (
	"encoding/json"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /seed", h.Seed)
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Seed(w http.ResponseWriter, r *http.Request) {
	result, ids, err := h.service.Seed(r.Context())
	if err != nil {
		apperror.HandleError(w, apperror.Internal("seed failed", err))
		return
	}

	response := map[string]interface{}{
		"message": "Seed completed successfully",
		"result":  result,
		"ids":     ids,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
