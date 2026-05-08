package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/testutil"
)

type mockRepo struct {
	createFunc       func(ctx context.Context, req *model.SupplierReturn) error
	getByIDFunc      func(ctx context.Context, id int) (*model.SupplierReturn, error)
	getAllFunc       func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error)
	updateStatusFunc func(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error)
	deleteFunc       func(ctx context.Context, id int) error
}

func (m *mockRepo) Create(ctx context.Context, req *model.SupplierReturn) error {
	return m.createFunc(ctx, req)
}
func (m *mockRepo) GetByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) GetAll(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
	return m.updateStatusFunc(ctx, id, status)
}
func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestCreateSupplierReturn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.CreateSupplierReturnMock()
		mock := &mockRepo{
			createFunc: func(ctx context.Context, req *model.SupplierReturn) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.CreateSupplierReturn(t.Context(), &expected)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("nil request returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		err := svc.CreateSupplierReturn(t.Context(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "request body is required" {
			t.Errorf("expected 'request body is required', got '%s'", err.Error())
		}
	})

	t.Run("empty return_no returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		req := testutil.CreateSupplierReturnMock()
		req.ReturnNo = ""
		err := svc.CreateSupplierReturn(t.Context(), &req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "return_no is required" {
			t.Errorf("expected 'return_no is required', got '%s'", err.Error())
		}
	})
}

func TestGetSupplierReturnByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.SupplierReturnMock()
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.SupplierReturn, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetSupplierReturnByID(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.SupplierReturnID != expected.SupplierReturnID {
			t.Errorf("expected ID %d, got %d", expected.SupplierReturnID, got.SupplierReturnID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.SupplierReturn, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetSupplierReturnByID(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGetAllSupplierReturns(t *testing.T) {
	t.Run("success with results", func(t *testing.T) {
		mock := &mockRepo{
		getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
			return []model.SupplierReturn{testutil.SupplierReturnMock()}, nil
		},
	}
	svc := service.NewService(mock)
	got, err := svc.GetAllSupplierReturns(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 supplier return, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &mockRepo{
		getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
			return []model.SupplierReturn{}, nil
		},
	}
	svc := service.NewService(mock)
	got, err := svc.GetAllSupplierReturns(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 supplier returns, got %d", len(got))
		}
	})

	t.Run("negative limit returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllSupplierReturns(t.Context(), model.QueryParams{Limit: -1, Offset: 0})
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
		_, err := svc.GetAllSupplierReturns(t.Context(), model.QueryParams{Limit: 10, Offset: -1})
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
		_, err := svc.GetAllSupplierReturns(t.Context(), model.QueryParams{Limit: 101, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})
}

func TestUpdateSupplierReturnStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.SupplierReturnMock()
		mock := &mockRepo{
			updateStatusFunc: func(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.UpdateSupplierReturnStatus(t.Context(), 1, model.ReturnStatusApproved)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.SupplierReturnID != expected.SupplierReturnID {
			t.Errorf("expected ID %d, got %d", expected.SupplierReturnID, got.SupplierReturnID)
		}
	})

	t.Run("invalid status returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.UpdateSupplierReturnStatus(t.Context(), 1, "invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "invalid supplier return status" {
			t.Errorf("expected 'invalid supplier return status', got '%s'", err.Error())
		}
	})
}

func TestDeleteSupplierReturn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteSupplierReturn(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return errors.New("supplier return not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteSupplierReturn(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "supplier return not found" {
			t.Errorf("expected 'supplier return not found', got '%s'", err.Error())
		}
	})
}
