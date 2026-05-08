package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/service"
	"github.com/google/uuid"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("GET /products", h.GetAll)
	mux.HandleFunc("GET /products/{id}", h.GetByID)
	mux.HandleFunc("PUT /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}

	if err := h.service.CreateProduct(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	result, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	params, err := parseGetProductsQueryParams(r)

	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	results, err := h.service.GetAllProducts(r.Context(), params)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	var req model.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.HandleError(w, apperror.BadRequest("invalid request body", err))
		return
	}
	req.ProductID = id

	if err := h.service.UpdateProduct(r.Context(), &req); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseProductID(r)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		apperror.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseProductID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperror.BadRequest("invalid product id", err)
	}
	if id == uuid.Nil {
		return uuid.Nil, apperror.BadRequest("invalid product id", nil)
	}
	return id, nil
}

func parseGetProductsQueryParams(r *http.Request) (model.GetProductsQueryParams, error) {
	q := r.URL.Query()

	var params model.GetProductsQueryParams

	category_id := q.Get("category")
	company_id := q.Get("company")

	if category_id != "" {
		catID, err := strconv.Atoi(category_id)
		if err != nil {
			return params, apperror.BadRequest("Invalid Category Id value", err)
		}
		params.CategoryID = &catID
	}

	if company_id != "" {
		compId, err := uuid.Parse(company_id)
		if err != nil {
			return params, apperror.BadRequest("Invalid Company Id value", err)
		}
		params.CompanyID = &compId
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return params, apperror.BadRequest("Invalid limit value", err)
		}
		params.Limit = limit
	}

	if offsetStr := q.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return params, apperror.BadRequest("Invalid offset value", err)
		}
		params.Offset = offset
	}

	return params, nil
}