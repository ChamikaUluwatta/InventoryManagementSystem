package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateInventory(ctx context.Context, inventory *model.Inventory) error {
	return s.repo.Create(ctx, inventory)
}

func (s *Service) GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllInventories(ctx context.Context) ([]model.Inventory, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateInventory(ctx context.Context, inventory *model.Inventory) error {
	if err := validation.ValidateStock(inventory.Stock); err != nil {
		return err
	}
	return s.repo.Update(ctx, inventory)
}

func (s *Service) DeleteInventory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetInventoriesByProduct(ctx context.Context, productID uuid.UUID) ([]model.Inventory, error) {
	return s.repo.GetByProduct(ctx, productID)
}

func (s *Service) GetInventoriesByLocation(ctx context.Context, locationID string) ([]model.Inventory, error) {
	return s.repo.GetByLocation(ctx, locationID)
}