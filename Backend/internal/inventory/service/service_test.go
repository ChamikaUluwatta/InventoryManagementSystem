package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/testutil"
	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc       func(ctx context.Context, inventory *model.Inventory) error
	getByIDFunc      func(ctx context.Context, id int) (*model.Inventory, error)
	getAllFunc       func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error)
	updateFunc       func(ctx context.Context, inventory *model.Inventory) error
	deleteFunc       func(ctx context.Context, id int) error
}

func (m *mockRepo) Create(ctx context.Context, inventory *model.Inventory) error {
	return m.createFunc(ctx, inventory)
}
func (m *mockRepo) GetByID(ctx context.Context, id int) (*model.Inventory, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) GetAll(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockRepo) Update(ctx context.Context, inventory *model.Inventory) error {
	return m.updateFunc(ctx, inventory)
}
func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestCreateInventory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, inv *model.Inventory) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		inv := testutil.CreateInventoryRequestMock()
		err := svc.CreateInventory(t.Context(), &inv)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("negative stock returns validation error", func(t *testing.T) {
		repoCalled := false
		mock := &mockRepo{
			createFunc: func(ctx context.Context, inv *model.Inventory) error {
				repoCalled = true
				return nil
			},
		}
		svc := service.NewService(mock)
		inv := testutil.CreateInventoryRequestMock()
		inv.Stock = -5
		err := svc.CreateInventory(t.Context(), &inv)
		if err == nil {
			t.Fatalf("expected validation error, got nil")
		}
		if err.Error() != "stock cannot be negative" {
			t.Errorf("expected 'stock cannot be negative', got '%s'", err.Error())
		}
		if repoCalled {
			t.Error("repo.Create should not be called when validation fails")
		}
	})
}

func TestGetInventoryByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.InventoryMock()
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.Inventory, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetInventoryByID(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.InventoryID != expected.InventoryID {
			t.Errorf("expected inventory ID %d, got %d", expected.InventoryID, got.InventoryID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.Inventory, error) {
				return nil, errors.New("inventory not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetInventoryByID(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}

func TestGetAllInventories(t *testing.T) {
	t.Run("success with results", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Inventory, error) {
				return []model.Inventory{testutil.InventoryMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllInventories(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 inventory, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, _ model.QueryParams) ([]model.Inventory, error) {
				return []model.Inventory{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllInventories(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 inventories, got %d", len(got))
		}
	})
	t.Run("negative limit returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllInventories(t.Context(), model.QueryParams{Limit: -1, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be non-negative" {
			t.Errorf("expected 'limit must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("negative offset returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllInventories(t.Context(), model.QueryParams{Limit: 10, Offset: -1})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "offset must be non-negative" {
			t.Errorf("expected 'offset must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("limit over max returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllInventories(t.Context(), model.QueryParams{Limit: 101, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})

	t.Run("passes product_id filter to repo", func(t *testing.T) {
		productID := uuid.New()
		var capturedParams model.QueryParams
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				capturedParams = params
				return []model.Inventory{}, nil
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetAllInventories(t.Context(), model.QueryParams{ProductID: &productID})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedParams.ProductID == nil || *capturedParams.ProductID != productID {
			t.Errorf("expected product_id %v, got %v", productID, capturedParams.ProductID)
		}
	})

	t.Run("passes location_id filter to repo", func(t *testing.T) {
		locationID := "LOC-001"
		var capturedParams model.QueryParams
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
				capturedParams = params
				return []model.Inventory{}, nil
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetAllInventories(t.Context(), model.QueryParams{LocationID: &locationID})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedParams.LocationID == nil || *capturedParams.LocationID != locationID {
			t.Errorf("expected location_id %s, got %v", locationID, capturedParams.LocationID)
		}
	})
}

func TestUpdateInventory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		inv := testutil.InventoryMock()
		err := svc.UpdateInventory(t.Context(), &inv)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("negative stock returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		inv := testutil.InventoryMock()
		inv.Stock = -1
		err := svc.UpdateInventory(t.Context(), &inv)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "stock cannot be negative" {
			t.Errorf("expected 'stock cannot be negative', got '%s'", err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, inv *model.Inventory) error {
				return errors.New("inventory not found")
			},
		}
		svc := service.NewService(mock)
		inv := testutil.InventoryMock()
		err := svc.UpdateInventory(t.Context(), &inv)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}

func TestDeleteInventory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteInventory(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return errors.New("inventory not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteInventory(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}


