package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidProductName = errors.New("invalid product name")
	ErrInvalidPrice       = errors.New("price cannot be negative")
)

type Service interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	GetProductsByCompany(ctx context.Context, companyID uuid.UUID) ([]Product, error)
	GetProductsByCategory(ctx context.Context, categoryID int) ([]Product, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return newService(repo)
}

func newService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, product *Product) error {
	if product.ProductName == "" {
		return ErrInvalidProductName
	}
	if product.Price.IsNegative() {
		return ErrInvalidPrice
	}
	return s.repo.Create(ctx, product)
}

func (s *service) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllProducts(ctx context.Context) ([]Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) UpdateProduct(ctx context.Context, product *Product) error {
	if product.ProductName == "" {
		return ErrInvalidProductName
	}
	if product.Price.IsNegative() {
		return ErrInvalidPrice
	}
	return s.repo.Update(ctx, product)
}

func (s *service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetProductsByCompany(ctx context.Context, companyID uuid.UUID) ([]Product, error) {
	return s.repo.GetByCompany(ctx, companyID)
}

func (s *service) GetProductsByCategory(ctx context.Context, categoryID int) ([]Product, error) {
	return s.repo.GetByCategory(ctx, categoryID)
}
