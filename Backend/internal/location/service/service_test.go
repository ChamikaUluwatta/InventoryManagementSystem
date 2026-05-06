package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/testutil"
)

type mockRepo struct {
	createFunc  func(ctx context.Context, location *model.Location) error
	getByIDFunc func(ctx context.Context, id string) (*model.Location, error)
	getAllFunc  func(ctx context.Context) ([]model.Location, error)
	updateFunc  func(ctx context.Context, location *model.Location) error
	deleteFunc  func(ctx context.Context, id string) error
}

func (m *mockRepo) Create(ctx context.Context, location *model.Location) error {
	return m.createFunc(ctx, location)
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.Location, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) GetAll(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
	return m.getAllFunc(ctx)
}
func (m *mockRepo) Update(ctx context.Context, location *model.Location) error {
	return m.updateFunc(ctx, location)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFunc(ctx, id)
}

func TestCreateLocation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		loc := testutil.LocationMock()
		err := svc.CreateLocation(t.Context(), &loc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty location id returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		loc := model.Location{LocationID: ""}
		err := svc.CreateLocation(t.Context(), &loc)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "location id is required" {
			t.Errorf("expected 'location id is required', got '%s'", err.Error())
		}
	})
}

func TestGetLocationByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.LocationMock()
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (*model.Location, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetLocationByID(t.Context(), expected.LocationID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.LocationID != expected.LocationID {
			t.Errorf("expected location ID '%s', got '%s'", expected.LocationID, got.LocationID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id string) (*model.Location, error) {
				return nil, errors.New("location not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetLocationByID(t.Context(), "NONEXIST")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}

func TestGetAllLocations(t *testing.T) {
	t.Run("success with results", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context) ([]model.Location, error) {
				return []model.Location{testutil.LocationMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllLocations(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 location, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context) ([]model.Location, error) {
				return []model.Location{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllLocations(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 locations, got %d", len(got))
		}
	})

}

func TestUpdateLocation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		loc := testutil.LocationMock()
		err := svc.UpdateLocation(t.Context(), &loc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty location id returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, l *model.Location) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		loc := model.Location{LocationID: ""}
		err := svc.UpdateLocation(t.Context(), &loc)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "location id is required" {
			t.Errorf("expected 'location id is required', got '%s'", err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, l *model.Location) error {
				return errors.New("location not found")
			},
		}
		svc := service.NewService(mock)
		loc := testutil.LocationMock()
		err := svc.UpdateLocation(t.Context(), &loc)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}

func TestDeleteLocation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteLocation(t.Context(), "TEST-LOC-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id string) error {
				return errors.New("location not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteLocation(t.Context(), "NONEXIST")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}
