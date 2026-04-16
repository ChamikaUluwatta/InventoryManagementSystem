package inventory

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /inventories", h.Create)
	mux.HandleFunc("GET /inventories", h.GetAll)
	mux.HandleFunc("GET /inventories/{id}", h.GetByID)
	mux.HandleFunc("PUT /inventories/{id}", h.Update)
	mux.HandleFunc("DELETE /inventories/{id}", h.Delete)
	mux.HandleFunc("GET /inventories/by-product/{productId}", h.GetByProduct)
	mux.HandleFunc("GET /inventories/by-location/{locationId}", h.GetByLocation)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req Inventory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.CreateInventory(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid inventory id", err))
		return
	}

	result, err := h.service.GetInventoryByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetAllInventories(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid inventory id", err))
		return
	}

	var req Inventory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}
	req.InventoryID = id

	if err := h.service.UpdateInventory(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid inventory id", err))
		return
	}

	if err := h.service.DeleteInventory(r.Context(), id); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetByProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("productId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid product id", err))
		return
	}

	results, err := h.service.GetInventoriesByProduct(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) GetByLocation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("locationId")
	if id == "" {
		apperror.HandleError(w, apperror.BadRequest("invalid location id", nil))
		return
	}

	results, err := h.service.GetInventoriesByLocation(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
