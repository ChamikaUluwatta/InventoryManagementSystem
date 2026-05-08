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
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
)

type mockService struct {
	createFunc    func(ctx context.Context, req *model.SupplierReturn) error
	getByIDFunc   func(ctx context.Context, id int) (*model.SupplierReturn, error)
	getAllFunc    func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error)
	updateFunc    func(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error)
	deleteFunc    func(ctx context.Context, id int) error
}

func (m *mockService) CreateSupplierReturn(ctx context.Context, req *model.SupplierReturn) error {
	return m.createFunc(ctx, req)
}
func (m *mockService) GetSupplierReturnByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllSupplierReturns(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateSupplierReturnStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
	return m.updateFunc(ctx, id, status)
}
func (m *mockService) DeleteSupplierReturn(ctx context.Context, id int) error {
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
			createFunc: func(ctx context.Context, req *model.SupplierReturn) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.CreateSupplierReturnMock())
		req := httptest.NewRequest("POST", "/supplier-returns", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/supplier-returns", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("validation error returns 400", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, req *model.SupplierReturn) error {
				return apperror.BadRequest("return_no is required", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.SupplierReturn{ReturnNo: ""})
		req := httptest.NewRequest("POST", "/supplier-returns", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	result := testutil.SupplierReturnWithItemsMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.SupplierReturn, error) {
				return &result, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.SupplierReturn, error) {
				return nil, apperror.NotFound("supplier return not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns/999", nil)
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
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
				return []model.SupplierReturn{testutil.SupplierReturnMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
				return []model.SupplierReturn{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("passes query params to service", func(t *testing.T) {
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
				capturedParams = params
				return []model.SupplierReturn{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/supplier-returns?limit=25&offset=10", nil)
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
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/supplier-returns", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestUpdateStatus(t *testing.T) {
	result := testutil.SupplierReturnMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
				return &result, nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.UpdateSupplierReturnStatusRequestMock())
		req := httptest.NewRequest("PATCH", "/supplier-returns/1/status", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.UpdateSupplierReturnStatusRequestMock())
		req := httptest.NewRequest("PATCH", "/supplier-returns/abc/status", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PATCH", "/supplier-returns/1/status", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid status returns 400", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
				return nil, apperror.BadRequest("invalid supplier return status", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.UpdateSupplierReturnStatusRequest{Status: "invalid"})
		req := httptest.NewRequest("PATCH", "/supplier-returns/1/status", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
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

		req := httptest.NewRequest("DELETE", "/supplier-returns/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/supplier-returns/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int) error {
				return apperror.NotFound("supplier return not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/supplier-returns/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, req *model.SupplierReturn) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.CreateSupplierReturnMock())
	req := httptest.NewRequest("POST", "/supplier-returns", body)
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
