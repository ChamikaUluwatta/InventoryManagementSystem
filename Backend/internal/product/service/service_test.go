package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/testutil"
	"github.com/google/uuid"
)

type mockRepo struct {
	createFunc func(ctx context.Context, product *model.CreateProductRequest) (model.Product, error)
	getById    func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error)
	getAll     func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error)
	update     func(ctx context.Context, product *model.Product) error
	delete     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockRepo) Create(ctx context.Context, product *model.CreateProductRequest) (model.Product, error) {
	return m.createFunc(ctx, product)
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
	return m.getById(ctx, id)
}

func (m *mockRepo) GetAll(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
	return m.getAll(ctx, params)
}

func (m *mockRepo) Update(ctx context.Context, product *model.Product) error {
	return m.update(ctx, product)
}

func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.delete(ctx, id)
}

const (
	invalidProductNameErrorMessage = "Invalid Product name"
	expectedErrorMessage           = "price cannot be negative"
	invalidCompanyIDErrorMessage   = "Invalid company id"
	invalidCategoryIDErrorMessage  = "Invalid category id"
)

func TestCreateProduct(t *testing.T) {
	mock := &mockRepo{
		createFunc: func(ctx context.Context, product *model.CreateProductRequest) (model.Product, error) {
			{
				return model.Product{}, nil
			}
		},
	}

	service := service.NewService(mock)

	t.Run("Validate Invalid Product Name", func(t *testing.T) {
		invalidProduct := testutil.CreateProductRequestMock()
		invalidProduct.ProductName = ""
		_, err := service.CreateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Error("Expected validation error for invalid product name but got nil")
		}

		if err != nil && err.Error() != invalidProductNameErrorMessage {
			t.Errorf("Expected error message '%s' but got '%s'", invalidProductNameErrorMessage, err.Error())
		}
	})

	t.Run("Validate Negative Price", func(t *testing.T) {
		invalidProduct := testutil.CreateProductRequestMock()
		invalidProduct.Price = invalidProduct.Price.Neg()
		_, err := service.CreateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Error("Expected validation error for negative price but got nil")
		}
		if err != nil && err.Error() != expectedErrorMessage {
			t.Errorf("Expected error message '%s' but got '%s'", expectedErrorMessage, err.Error())
		}
	})

	t.Run("Validate Invalid Diameter and Width", func(t *testing.T) {
		invalidProductDiameter := testutil.CreateProductRequestMock()
		invalidProductDiameter.Diameter = invalidProductDiameter.Diameter.Neg()
		invalidProductDiameter.Width = invalidProductDiameter.Width.Neg()
		_, err := service.CreateProduct(t.Context(), &invalidProductDiameter)
		if err == nil {
			t.Fatalf("Expected validation error for negative diameter and width but got nil")
		}

		if err.Error() != "diameter cannot be negative" {
			t.Errorf("Expected error message 'diameter cannot be negative' but got '%s'", err.Error())
		}

		invalidProductWidth := testutil.CreateProductRequestMock()
		invalidProductWidth.Width = invalidProductWidth.Width.Neg()
		_, err = service.CreateProduct(t.Context(), &invalidProductWidth)
		if err == nil {
			t.Fatalf("Expected validation error for negative width but got nil")
		}

		if err.Error() != "width cannot be negative" {
			t.Errorf("Expected error message 'width cannot be negative' but got '%s'", err.Error())
		}
	})

	t.Run("Validate Invalid Company ID", func(t *testing.T) {
		invalidProduct := testutil.CreateProductRequestMock()
		invalidProduct.CompanyID = uuid.Nil
		_, err := service.CreateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Error("Expected validation error for invalid company id but got nil")
		}
		if err != nil && err.Error() != invalidCompanyIDErrorMessage {
			t.Errorf("Expected error message '%s' but got '%s'", invalidCompanyIDErrorMessage, err.Error())
		}
	})

	t.Run("Validate Invalid Category ID", func(t *testing.T) {
		invalidProduct := testutil.CreateProductRequestMock()
		invalidProduct.CategoryID = 0
		_, err := service.CreateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Error("Expected validation error for invalid category id but got nil")
		}
		if err != nil && err.Error() != invalidCategoryIDErrorMessage {
			t.Errorf("Expected error message '%s' but got '%s'", invalidCategoryIDErrorMessage, err.Error())
		}
	})
}

func TestGetProductByID(t *testing.T) {
	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		productMock := testutil.GetProductByIdMock()
		mock := &mockRepo{
			getById: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				return &productMock, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetProductByID(t.Context(), productID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if got.ProductID != productMock.ProductID {
			t.Errorf("Expected product ID %v, got %v", productMock.ProductID, got.ProductID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getById: func(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
				return nil, errors.New("product not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetProductByID(t.Context(), productID)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("Expected 'product not found', got '%s'", err.Error())
		}
	})

}

func TestGetAllProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			getAll: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				return []model.Product{testutil.ProductMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllProducts(t.Context(), model.GetProductsQueryParams{})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("Expected 1 product, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var paramsCapture model.GetProductsQueryParams
		mock := &mockRepo{
			getAll: func(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
				paramsCapture = params
				return []model.Product{}, nil
			},
		}

		if paramsCapture.CategoryID != nil || paramsCapture.CompanyID != nil {
			t.Errorf("Expected empty query params, got %v", paramsCapture)
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllProducts(t.Context(), model.GetProductsQueryParams{})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("Expected 0 products, got %d", len(got))
		}
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			update: func(ctx context.Context, product *model.Product) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		p := testutil.ProductMock()
		err := svc.UpdateProduct(t.Context(), &p)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("invalid product name", func(t *testing.T) {
		mock := &mockRepo{
			update: func(ctx context.Context, product *model.Product) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		invalidProduct := testutil.ProductMock()
		invalidProduct.ProductName = ""
		err := svc.UpdateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if err.Error() != invalidProductNameErrorMessage {
			t.Errorf("Expected '%s', got '%s'", invalidProductNameErrorMessage, err.Error())
		}
	})

	t.Run("negative price", func(t *testing.T) {
		mock := &mockRepo{
			update: func(ctx context.Context, product *model.Product) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		invalidProduct := testutil.ProductMock()
		invalidProduct.Price = invalidProduct.Price.Neg()
		err := svc.UpdateProduct(t.Context(), &invalidProduct)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
		if err.Error() != expectedErrorMessage {
			t.Errorf("Expected '%s', got '%s'", expectedErrorMessage, err.Error())
		}
	})
}

func TestDeleteProduct(t *testing.T) {
	productID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			delete: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteProduct(t.Context(), productID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			delete: func(ctx context.Context, id uuid.UUID) error {
				return errors.New("product not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteProduct(t.Context(), productID)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("Expected 'product not found', got '%s'", err.Error())
		}
	})
}
