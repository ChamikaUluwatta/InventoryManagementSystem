package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

type mockService struct {
	createFunc  func(ctx context.Context, product *model.Product) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error)
	getAllFunc  func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error)
	updateFunc  func(ctx context.Context, product *model.Product) error
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockService) CreateProduct(ctx context.Context, product *model.Product) error {
	return m.createFunc(ctx, product)
}
func (m *mockService) GetProductByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllProducts(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateProduct(ctx context.Context, product *model.Product) error {
	return m.updateFunc(ctx, product)
}
func (m *mockService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return m.deleteFunc(ctx, id)
}

func setupHandler(svc service.Service) *http.ServeMux {
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}


func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, p *model.Product) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.ProductMock())
		req := httptest.NewRequest("POST", "/products", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/products", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, p *model.Product) error {
				return apperror.BadRequest("Invalid Product name", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.ProductMock())
		req := httptest.NewRequest("POST", "/products", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	product := testutil.ProductMock()

	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				return &model.GetProductById{Product: product, Stock: 10}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products/"+product.ProductID.String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products/not-a-uuid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("nil UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products/00000000-0000-0000-0000-000000000000", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				return nil, apperror.NotFound("product not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				return []model.Product{testutil.ProductMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("with filters", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				if params.CategoryID == nil || params.CompanyID == nil {
					t.Error("expected category and company filters")
				}
				return []model.Product{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products?category=1&company="+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid category filter", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products?category=abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid company filter", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products?company=not-a-uuid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/products", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	product := testutil.ProductMock()

	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, p *model.Product) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, product)
		req := httptest.NewRequest("PUT", "/products/"+product.ProductID.String(), body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/products/not-a-uuid", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("nil UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/products/00000000-0000-0000-0000-000000000000", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/products/"+product.ProductID.String(), strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, p *model.Product) error {
				return apperror.NotFound("product not found", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, product)
		req := httptest.NewRequest("PUT", "/products/"+uuid.New().String(), body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/products/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/products/not-a-uuid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("nil UUID", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/products/00000000-0000-0000-0000-000000000000", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return apperror.NotFound("product not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/products/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, p *model.Product) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.ProductMock())
	req := httptest.NewRequest("POST", "/products", body)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}

	var resp apperror.AppError
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code 500 in body, got %d", resp.StatusCode)
	}
}
