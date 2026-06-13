package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/auth"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/companies", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(auth.RequirePermission(PermissionWrite))
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
		r.Group(func(r chi.Router) {
			r.Use(auth.RequirePermission(PermissionRead))
			r.Get("/", h.GetAll)
			r.Get("/{id}", h.GetByID)
			r.Get("/{id}/dependencies", h.GetCompanyDependencies)
		})
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Company
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.CreateCompany(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid company id", err))
		return
	}

	result, err := h.service.GetCompanyByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	params := parseQueryParams(r)
	results, err := h.service.GetAllCompanies(r.Context(), params)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid company id", err))
		return
	}

	var req model.Company
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}
	req.CompanyID = id

	if err := h.service.UpdateCompany(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid company id", err))
		return
	}

	if err := h.service.DeleteCompany(r.Context(), id); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCompanyDependencies(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid company id", err))
		return
	}

	result, err := h.service.GetCompanyDependencies(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
