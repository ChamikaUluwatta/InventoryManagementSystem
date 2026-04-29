package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /categories", h.Create)
	mux.HandleFunc("GET /categories", h.GetAll)
	mux.HandleFunc("GET /categories/{id}", h.GetByID)
	mux.HandleFunc("PUT /categories/{id}", h.Update)
	mux.HandleFunc("DELETE /categories/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.CreateCategory(r.Context(), &req); err != nil {
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
		apperror.HandleError(w, apperror.BadRequest("invalid category id", err))
		return
	}

	result, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetAllCategories(r.Context())
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
		apperror.HandleError(w, apperror.BadRequest("invalid category id", err))
		return
	}

	var req model.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}
	req.CategoryID = id

	if err := h.service.UpdateCategory(r.Context(), &req); err != nil {
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
		apperror.HandleError(w, apperror.BadRequest("invalid category id", err))
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}