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
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

type mockService struct {
	createFunc  func(ctx context.Context, inventory *model.Inventory) error
	getByIDFunc func(ctx context.Context, id int) (*model.Inventory, error)
	getAllFunc  func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error)
	updateFunc  func(ctx context.Context, inventory *model.Inventory) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *mockService) CreateInventory(ctx context.Context, inventory *model.Inventory) error {
	return m.createFunc(ctx, inventory)
}
func (m *mockService) GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllInventories(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateInventory(ctx context.Context, inventory *model.Inventory) error {
	return m.updateFunc(ctx, inventory)
}
func (m *mockService) DeleteInventory(ctx context.Context, id int) error {
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
			createFunc: func(ctx context.Context, inv *model.Inventory) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.InventoryMock())
		req := httptest.NewRequest("POST", "/inventories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/inventories", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("FK violation returns 400", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, inv *model.Inventory) error {
				return apperror.BadRequest("invalid product_id or location_id", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.InventoryMock())
		req := httptest.NewRequest("POST", "/inventories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("unique violation returns 409", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, inv *model.Inventory) error {
				return apperror.Conflict("inventory already exists for product and location", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.InventoryMock())
		req := httptest.NewRequest("POST", "/inventories", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	inventory := testutil.InventoryMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.Inventory, error) {
				return &inventory, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int) (*model.Inventory, error) {
				return nil, apperror.NotFound("inventory not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories/999", nil)
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
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				return []model.Inventory{testutil.InventoryMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				return []model.Inventory{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/inventories", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
	t.Run("passes query params to service", func(t *testing.T) {
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				capturedParams = params
				return []model.Inventory{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/inventories?limit=25&offset=10", nil)
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

	t.Run("filters by product_id", func(t *testing.T) {
		productID := uuid.New()
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				capturedParams = params
				return []model.Inventory{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/inventories?product_id="+productID.String(), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if capturedParams.ProductID == nil || *capturedParams.ProductID != productID {
			t.Errorf("expected product_id %v, got %v", productID, capturedParams.ProductID)
		}
	})

	t.Run("filters by location_id", func(t *testing.T) {
		locationID := "LOC-001"
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				capturedParams = params
				return []model.Inventory{}, nil
			},
		}
		mux := setupHandler(svc)
		req := httptest.NewRequest("GET", "/inventories?location_id="+locationID, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
		if capturedParams.LocationID == nil || *capturedParams.LocationID != locationID {
			t.Errorf("expected location_id %s, got %v", locationID, capturedParams.LocationID)
		}
	})
}

func TestUpdate(t *testing.T) {
	inventory := testutil.InventoryMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, inventory)
		req := httptest.NewRequest("PUT", "/inventories/1", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/inventories/abc", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/inventories/1", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("negative stock returns 400", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return apperror.BadRequest("stock cannot be negative", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.Inventory{Stock: -1})
		req := httptest.NewRequest("PUT", "/inventories/1", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return apperror.NotFound("inventory not found", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, inventory)
		req := httptest.NewRequest("PUT", "/inventories/999", body)
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

		req := httptest.NewRequest("DELETE", "/inventories/1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/inventories/abc", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int) error {
				return apperror.NotFound("inventory not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/inventories/999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, inv *model.Inventory) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.InventoryMock())
	req := httptest.NewRequest("POST", "/inventories", body)
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
