package inventory

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateInventory(ctx context.Context, inventory *Inventory) error {
	if inventory.Stock < 0 {
		return apperror.BadRequest("stock cannot be negative", nil)
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
		return apperror.BadRequest("stock cannot be negative", nil)
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
