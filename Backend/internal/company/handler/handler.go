package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/service"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /companies", h.Create)
	mux.HandleFunc("GET /companies", h.GetAll)
	mux.HandleFunc("GET /companies/{id}", h.GetByID)
	mux.HandleFunc("PUT /companies/{id}", h.Update)
	mux.HandleFunc("DELETE /companies/{id}", h.Delete)
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
	idStr := r.PathValue("id")
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
	results, err := h.service.GetAllCompanies(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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
	idStr := r.PathValue("id")
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