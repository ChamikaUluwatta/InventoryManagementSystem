package product

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*GetProductById, error)
	GetAllProducts(ctx context.Context, params GetProductsQueryParams) ([]Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, product *Product) error {
	if product.ProductName == "" {
		return apperror.BadRequest("product name is required", nil)
	}
	if product.Price.IsNegative() {
		return apperror.BadRequest("price cannot be negative", nil)
	}
	return s.repo.Create(ctx, product)
}

func (s *service) GetProductByID(ctx context.Context, id uuid.UUID) (*GetProductById, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllProducts(ctx context.Context, params GetProductsQueryParams) ([]Product, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *service) UpdateProduct(ctx context.Context, product *Product) error {
	if product.ProductName == "" {
		return apperror.BadRequest("product name is required", nil)
	}
	if product.Price.IsNegative() {
		return apperror.BadRequest("price cannot be negative", nil)
	}
	return s.repo.Update(ctx, product)
}

func (s *service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
