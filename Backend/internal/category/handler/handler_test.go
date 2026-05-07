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
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
)

type mockService struct {
	createFunc  func(ctx context.Context, category *model.Category) error
	getByIDFunc func(ctx context.Context, id int) (*model.Category, error)
	getAllFunc  func(ctx context.Context, params model.QueryParams) ([]model.Category, error)
	updateFunc  func(ctx context.Context, category *model.Category) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockService) CreateCategory(ctx context.Context, category *model.Category) error {
	return m.createFunc(ctx, category)
}
func (m *mockService) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllCategories(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateCategory(ctx context.Context, category *model.Category) error {
	return m.updateFunc(ctx, category)
}
func (m *mockService) DeleteCategory(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}
func (m *mockService) GetCategoriesByParent(ctx context.Context, parentID *int) ([]model.Category, error) {
	return nil, nil
}

func setupHandler(svc service.Service) *http.ServeMux {
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}


func TestCreate(t *testing.T) {
	t.Run("success returns 201", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.CategoryMock())
		req := httptest.NewRequest("POST", "/categories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/categories", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service error returns 400", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return apperror.BadRequest("category name is required", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.Category{CategoryName: ""})
		req := httptest.NewRequest("POST", "/categories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service Internal error returns 500", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.CategoryMock())
		req := httptest.NewRequest("POST", "/categories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	category := testutil.CategoryMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.Category, error) {
				return &category, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.Category, error) {
				return nil, apperror.NotFound("category not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("success with results returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				return []model.Category{testutil.CategoryMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				return []model.Category{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("passes query params to service", func(t *testing.T) {
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				capturedParams = params
				return []model.Category{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/categories?limit=25&offset=10", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if capturedParams.Limit != 25 {
			t.Errorf("expected limit 25, got %d", capturedParams.Limit)
		}
		if capturedParams.Offset != 10 {
			t.Errorf("expected offset 10, got %d", capturedParams.Offset)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/categories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	category := testutil.CategoryMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, category)
		req := httptest.NewRequest("PUT", "/categories/1", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/categories/abc", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/categories/1", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, c *model.Category) error {
				return apperror.NotFound("category not found", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, category)
		req := httptest.NewRequest("PUT", "/categories/999", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success returns 204", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/categories/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/categories/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int) error {
				return apperror.NotFound("category not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/categories/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, c *model.Category) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.CategoryMock())
	req := httptest.NewRequest("POST", "/categories", body)
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
