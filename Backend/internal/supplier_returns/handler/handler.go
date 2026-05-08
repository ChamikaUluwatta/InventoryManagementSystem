package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /supplier-returns", h.Create)
	mux.HandleFunc("GET /supplier-returns", h.GetAll)
	mux.HandleFunc("GET /supplier-returns/{id}", h.GetByID)
	mux.HandleFunc("PATCH /supplier-returns/{id}/status", h.UpdateStatus)
	mux.HandleFunc("DELETE /supplier-returns/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var result model.SupplierReturn
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	err := h.service.CreateSupplierReturn(r.Context(), &result)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid supplier return id", err))
		return
	}

	result, err := h.service.GetSupplierReturnByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := parseQueryParams(r)
	results, err := h.service.GetAllSupplierReturns(r.Context(), params)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid supplier return id", err))
		return
	}

	var req model.UpdateSupplierReturnStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	result, err := h.service.UpdateSupplierReturnStatus(r.Context(), id, req.Status)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid supplier return id", err))
		return
	}

	if err := h.service.DeleteSupplierReturn(r.Context(), id); err != nil {
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
	return params
}
