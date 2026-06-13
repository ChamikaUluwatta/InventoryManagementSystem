package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/validation"
)

type Service interface {
	CreateInventory(ctx context.Context, inventory *model.Inventory) error
	GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error)
	GetAllInventories(ctx context.Context, params model.QueryParams) ([]model.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *model.Inventory) error
	DeleteInventory(ctx context.Context, id int) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateInventory(ctx context.Context, inventory *model.Inventory) error {
	if err := validation.ValidateStock(inventory.Stock); err != nil {
		return err
	}
	return s.repo.Create(ctx, inventory)
}

func (s *service) GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllInventories(ctx context.Context, params model.QueryParams) ([]model.Inventory, error) {
	validatedParams, err := validation.ValidateParams(params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	validatedParams.ProductID = params.ProductID
	validatedParams.LocationID = params.LocationID
	return s.repo.GetAll(ctx, validatedParams)
}

func (s *service) UpdateInventory(ctx context.Context, inventory *model.Inventory) error {
	if err := validation.ValidateStock(inventory.Stock); err != nil {
		return err
	}
	return s.repo.Update(ctx, inventory)
}

func (s *service) DeleteInventory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
