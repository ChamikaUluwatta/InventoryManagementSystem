package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/validation"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLocation(ctx context.Context, location *model.Location) error {
	if err := validation.ValidateLocationID(location.LocationID); err != nil {
		return err
	}
	return s.repo.Create(ctx, location)
}

func (s *Service) GetLocationByID(ctx context.Context, id string) (*model.Location, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllLocations(ctx context.Context) ([]model.Location, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateLocation(ctx context.Context, location *model.Location) error {
	if err := validation.ValidateLocationID(location.LocationID); err != nil {
		return err
	}
	return s.repo.Update(ctx, location)
}

func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}