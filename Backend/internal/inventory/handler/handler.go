package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/service"
	"github.com/google/uuid"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /inventories", h.Create)
	mux.HandleFunc("GET /inventories", h.GetAll)
	mux.HandleFunc("GET /inventories/{id}", h.GetByID)
	mux.HandleFunc("PUT /inventories/{id}", h.Update)
	mux.HandleFunc("DELETE /inventories/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Inventory
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
	params := parseQueryParams(r)
	results, err := h.service.GetAllInventories(r.Context(), params)
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

	var req model.Inventory
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

func parseQueryParams(r *http.Request) model.QueryParams {
	var params model.QueryParams

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}
	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		if id, err := uuid.Parse(productIDStr); err == nil {
			params.ProductID = &id
		}
	}
	if locationIDStr := r.URL.Query().Get("location_id"); locationIDStr != "" {
		params.LocationID = &locationIDStr
	}

	return params
}