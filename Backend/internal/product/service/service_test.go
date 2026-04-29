package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type mockRepository struct {
	CreateFunc    func(ctx context.Context, product *model.Product) error
	GetByIDFunc   func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error)
	GetAllFunc    func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error)
	UpdateFunc    func(ctx context.Context, product *model.Product) error
	DeleteFunc    func(ctx context.Context, id uuid.UUID) error
}

func (m *mockRepository) Create(ctx context.Context, product *model.Product) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, product)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetAll(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, params)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, product *model.Product) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, product)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		product *model.Product
		mockRepo *mockRepository
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid product",
			product: &model.Product{
				ProductName: "Test Product",
				Price:       decimal.NewFromFloat(99.99),
			},
			mockRepo: &mockRepository{
				CreateFunc: func(ctx context.Context, product *model.Product) error {
					product.ProductID = uuid.New()
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "empty product name",
			product: &model.Product{
				ProductName: "",
				Price:       decimal.NewFromFloat(99.99),
			},
			wantErr: true,
			errMsg:  "product name is required",
		},
		{
			name: "negative price",
			product: &model.Product{
				ProductName: "Test Product",
				Price:       decimal.NewFromFloat(-10.00),
			},
			wantErr: true,
			errMsg:  "price cannot be negative",
		},
		{
			name: "zero price is valid",
			product: &model.Product{
				ProductName: "Free Product",
				Price:       decimal.Zero,
			},
			mockRepo: &mockRepository{
				CreateFunc: func(ctx context.Context, product *model.Product) error {
					product.ProductID = uuid.New()
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.mockRepo)
			err := svc.CreateProduct(context.Background(), tt.product)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				var appErr *apperror.AppError
				if errors.As(err, &appErr) {
					if appErr.Message != tt.errMsg {
						t.Errorf("expected error message %q, got %q", tt.errMsg, appErr.Message)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name    string
		product *model.Product
		mockRepo *mockRepository
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update",
			product: &model.Product{
				ProductID:   uuid.New(),
				ProductName: "Updated Product",
				Price:       decimal.NewFromFloat(149.99),
			},
			mockRepo: &mockRepository{
				UpdateFunc: func(ctx context.Context, product *model.Product) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "empty product name",
			product: &model.Product{
				ProductID:   uuid.New(),
				ProductName: "",
				Price:       decimal.NewFromFloat(149.99),
			},
			wantErr: true,
			errMsg:  "product name is required",
		},
		{
			name: "negative price",
			product: &model.Product{
				ProductID:   uuid.New(),
				ProductName: "Test Product",
				Price:       decimal.NewFromFloat(-50.00),
			},
			wantErr: true,
			errMsg:  "price cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.mockRepo)
			err := svc.UpdateProduct(context.Background(), tt.product)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				var appErr *apperror.AppError
				if errors.As(err, &appErr) {
					if appErr.Message != tt.errMsg {
						t.Errorf("expected error message %q, got %q", tt.errMsg, appErr.Message)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	t.Run("product found", func(t *testing.T) {
		expectedID := uuid.New()
		mockRepo := &mockRepository{
			GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				if id != expectedID {
					t.Errorf("expected id %v, got %v", expectedID, id)
				}
				return &model.GetProductById{
					Product: model.Product{
						ProductID:   expectedID,
						ProductName: "Test Product",
						Price:       decimal.NewFromFloat(99.99),
					},
					Stock: 100,
				}, nil
			},
		}

		svc := NewService(mockRepo)
		result, err := svc.GetProductByID(context.Background(), expectedID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected result but got nil")
		}
		if result.Stock != 100 {
			t.Errorf("expected stock 100, got %d", result.Stock)
		}
	})

	t.Run("product not found", func(t *testing.T) {
		notFoundID := uuid.New()
		mockRepo := &mockRepository{
			GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				return nil, apperror.NotFound("product not found", nil)
			},
		}

		svc := NewService(mockRepo)
		result, err := svc.GetProductByID(context.Background(), notFoundID)

		if err == nil {
			t.Error("expected error but got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("returns all products", func(t *testing.T) {
		expectedProducts := []model.Product{
			{ProductID: uuid.New(), ProductName: "Product A", Price: decimal.NewFromFloat(10.00)},
			{ProductID: uuid.New(), ProductName: "Product B", Price: decimal.NewFromFloat(20.00)},
		}

		mockRepo := &mockRepository{
			GetAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				return expectedProducts, nil
			},
		}

		svc := NewService(mockRepo)
		results, err := svc.GetAllProducts(context.Background(), model.GetProductsQueryParams{})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("expected 2 products, got %d", len(results))
		}
	})

	t.Run("filters by category", func(t *testing.T) {
		catID := 1
		mockRepo := &mockRepository{
			GetAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				if params.CategoryID == nil || *params.CategoryID != catID {
					t.Errorf("expected category filter %d, got %v", catID, params.CategoryID)
				}
				return []model.Product{}, nil
			},
		}

		svc := NewService(mockRepo)
		_, err := svc.GetAllProducts(context.Background(), model.GetProductsQueryParams{CategoryID: &catID})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("filters by company", func(t *testing.T) {
		companyID := uuid.New()
		mockRepo := &mockRepository{
			GetAllFunc: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				if params.CompanyID == nil || *params.CompanyID != companyID {
					t.Errorf("expected company filter %v, got %v", companyID, params.CompanyID)
				}
				return []model.Product{}, nil
			},
		}

		svc := NewService(mockRepo)
		_, err := svc.GetAllProducts(context.Background(), model.GetProductsQueryParams{CompanyID: &companyID})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("delete succeeds", func(t *testing.T) {
		deleteID := uuid.New()
		called := false
		mockRepo := &mockRepository{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				called = true
				if id != deleteID {
					t.Errorf("expected id %v, got %v", deleteID, id)
				}
				return nil
			},
		}

		svc := NewService(mockRepo)
		err := svc.DeleteProduct(context.Background(), deleteID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("delete function was not called")
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		deleteID := uuid.New()
		mockRepo := &mockRepository{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return apperror.NotFound("product not found", nil)
			},
		}

		svc := NewService(mockRepo)
		err := svc.DeleteProduct(context.Background(), deleteID)

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}