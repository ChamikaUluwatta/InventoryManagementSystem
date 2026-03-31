package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidStock = errors.New("stock cannot be negative")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateInventory(ctx context.Context, inventory *Inventory) error {
	if inventory.Stock < 0 {
		return ErrInvalidStock
	}
	return s.repo.Create(ctx, inventory)
}

func (s *Service) GetInventoryByID(ctx context.Context, id int) (*Inventory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllInventories(ctx context.Context) ([]Inventory, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateInventory(ctx context.Context, inventory *Inventory) error {
	if inventory.Stock < 0 {
		return ErrInvalidStock
	}
	return s.repo.Update(ctx, inventory)
}

func (s *Service) DeleteInventory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetInventoriesByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error) {
	return s.repo.GetByProduct(ctx, productID)
}

func (s *Service) GetInventoriesByLocation(ctx context.Context, locationID string) ([]Inventory, error) {
	return s.repo.GetByLocation(ctx, locationID)
}
