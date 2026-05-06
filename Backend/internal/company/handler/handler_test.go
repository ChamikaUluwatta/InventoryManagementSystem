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
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

type mockService struct {
	createFunc  func(ctx context.Context, company *model.Company) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*model.Company, error)
	getAllFunc  func(ctx context.Context, params model.QueryParams) ([]model.Company, error)
	updateFunc  func(ctx context.Context, company *model.Company) error
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockService) CreateCompany(ctx context.Context, company *model.Company) error {
	return m.createFunc(ctx, company)
}
func (m *mockService) GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllCompanies(ctx context.Context, params model.QueryParams) ([]model.Company, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateCompany(ctx context.Context, company *model.Company) error {
	return m.updateFunc(ctx, company)
}
func (m *mockService) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return m.deleteFunc(ctx, id)
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
			createFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.CompanyMock())
		req := httptest.NewRequest("POST", "/companies", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/companies", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service error returns 400", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, c *model.Company) error {
				return apperror.BadRequest("company name is required", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.Company{CompanyName: ""})
		req := httptest.NewRequest("POST", "/companies", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service Internal error returns 500", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, c *model.Company) error {
				return apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.CompanyMock())
		req := httptest.NewRequest("POST", "/companies", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	company := testutil.CompanyMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Company, error) {
				return &company, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies/"+company.CompanyID.String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies/not-a-uuid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.Company, error) {
				return nil, apperror.NotFound("company not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies/"+uuid.New().String(), nil)
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
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Company, error) {
				return []model.Company{testutil.CompanyMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Company, error) {
				return []model.Company{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("passes query params to service", func(t *testing.T) {
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Company, error) {
				capturedParams = params
				return []model.Company{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/companies?limit=25&offset=10", nil)
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
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Company, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/companies", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	company := testutil.CompanyMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, c *model.Company) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, company)
		req := httptest.NewRequest("PUT", "/companies/"+company.CompanyID.String(), body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/companies/not-a-uuid", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/companies/"+company.CompanyID.String(), strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, c *model.Company) error {
				return apperror.NotFound("company not found", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, company)
		req := httptest.NewRequest("PUT", "/companies/"+uuid.New().String(), body)
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
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/companies/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("invalid UUID returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/companies/not-a-uuid", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return apperror.NotFound("company not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/companies/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, c *model.Company) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.CompanyMock())
	req := httptest.NewRequest("POST", "/companies", body)
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
