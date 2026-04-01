package location

import (
	"context"
	"errors"
)

var ErrInvalidLocationID = errors.New("invalid location id")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLocation(ctx context.Context, location *Location) error {
	if location.LocationID == "" {
		return ErrInvalidLocationID
	}
	return s.repo.Create(ctx, location)
}

func (s *Service) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) GetAllLocations(ctx context.Context) ([]Location, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateLocation(ctx context.Context, location *Location) error {
	if location.LocationID == "" {
		return ErrInvalidLocationID
	}
	return s.repo.Update(ctx, location)
}

func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
