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
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/handler"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/testutil"
	sharedtestutil "github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
)

type mockService struct {
	createFunc  func(ctx context.Context, location *model.Location) error
	getByIDFunc func(ctx context.Context, id string) (*model.Location, error)
	getAllFunc  func(ctx context.Context, params model.QueryParams) ([]model.Location, error)
	updateFunc  func(ctx context.Context, location *model.Location) error
	deleteFunc  func(ctx context.Context, id string) error
}

func (m *mockService) CreateLocation(ctx context.Context, location *model.Location) error {
	return m.createFunc(ctx, location)
}
func (m *mockService) GetLocationByID(ctx context.Context, id string) (*model.Location, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockService) GetAllLocations(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockService) UpdateLocation(ctx context.Context, location *model.Location) error {
	return m.updateFunc(ctx, location)
}
func (m *mockService) DeleteLocation(ctx context.Context, id string) error {
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
			createFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.LocationMock())
		req := httptest.NewRequest("POST", "/locations", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("POST", "/locations", strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service error returns 400", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, l *model.Location) error {
				return apperror.BadRequest("location id is required", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, model.Location{LocationID: ""})
		req := httptest.NewRequest("POST", "/locations", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("service Internal error returns 500", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, l *model.Location) error {
				return apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, testutil.LocationMock())
		req := httptest.NewRequest("POST", "/locations", body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestGetByID(t *testing.T) {
	location := testutil.LocationMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id string) (*model.Location, error) {
				return &location, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations/"+location.LocationID, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id string) (*model.Location, error) {
				return nil, apperror.NotFound("location not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations/NONEXIST", nil)
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
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
				return []model.Location{testutil.LocationMock()}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
				return []model.Location{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("passes query params to service", func(t *testing.T) {
		var capturedParams model.QueryParams
		svc := &mockService{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
				capturedParams = params
				return []model.Location{}, nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations?limit=25&offset=10", nil)
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
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
				return nil, apperror.Internal("db error", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("GET", "/locations", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rec.Code)
		}
	})
}

func TestUpdate(t *testing.T) {
	location := testutil.LocationMock()

	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, location)
		req := httptest.NewRequest("PUT", "/locations/"+location.LocationID, body)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		svc := &mockService{}
		mux := setupHandler(svc)

		req := httptest.NewRequest("PUT", "/locations/"+location.LocationID, strings.NewReader("{invalid"))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, l *model.Location) error {
				return apperror.NotFound("location not found", nil)
			},
		}
		mux := setupHandler(svc)

		body := sharedtestutil.MarshalBody(t, location)
		req := httptest.NewRequest("PUT", "/locations/NONEXIST", body)
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
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/locations/TEST-LOC-1", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rec.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id string) error {
				return apperror.NotFound("location not found", nil)
			},
		}
		mux := setupHandler(svc)

		req := httptest.NewRequest("DELETE", "/locations/NONEXIST", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestErrorResponse(t *testing.T) {
	svc := &mockService{
		createFunc: func(ctx context.Context, l *model.Location) error {
			return errors.New("unexpected raw error")
		},
	}
	mux := setupHandler(svc)

	body := sharedtestutil.MarshalBody(t, testutil.LocationMock())
	req := httptest.NewRequest("POST", "/locations", body)
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
