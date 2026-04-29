package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/validation"
	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error)
	GetAllProducts(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateProduct(ctx context.Context, product *model.Product) error {
	if err := validation.ValidateProductName(product.ProductName); err != nil {
		return err
	}
	if err := validation.ValidatePrice(product.Price); err != nil {
		return err
	}
	return s.repo.Create(ctx, product)
}

func (s *service) GetProductByID(ctx context.Context, id uuid.UUID) (*model.GetProductById, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllProducts(ctx context.Context, params model.GetProductsQueryParams) ([]model.Product, error) {
	return s.repo.GetAll(ctx, params)
}

func (s *service) UpdateProduct(ctx context.Context, product *model.Product) error {
	if err := validation.ValidateProductName(product.ProductName); err != nil {
		return err
	}
	if err := validation.ValidatePrice(product.Price); err != nil {
		return err
	}
	return s.repo.Update(ctx, product)
}

func (s *service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}