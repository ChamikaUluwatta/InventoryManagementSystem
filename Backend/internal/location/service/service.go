package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/validation"
)

type Service interface {
	CreateLocation(ctx context.Context, location *model.Location) error
	GetLocationByID(ctx context.Context, id string) (*model.Location, error)
	GetAllLocations(ctx context.Context, params model.QueryParams) ([]model.Location, error)
	UpdateLocation(ctx context.Context, location *model.Location) error
	DeleteLocation(ctx context.Context, id string) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateLocation(ctx context.Context, location *model.Location) error {
	if err := validation.ValidateLocationID(location.LocationID); err != nil {
		return err
	}
	return s.repo.Create(ctx, location)
}

func (s *service) GetLocationByID(ctx context.Context, id string) (*model.Location, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllLocations(ctx context.Context, params model.QueryParams) ([]model.Location, error) {
	if err, validatedParams := validation.ValidateParams(params.Limit, params.Offset); err != nil {
		return nil, err
	} else {
		params = validatedParams
	}
	return s.repo.GetAll(ctx, params)
}

func (s *service) UpdateLocation(ctx context.Context, location *model.Location) error {
	if err := validation.ValidateLocationID(location.LocationID); err != nil {
		return err
	}
	return s.repo.Update(ctx, location)
}

func (s *service) DeleteLocation(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
