package company

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

func (s *Service) CreateCompany(ctx context.Context, company *Company) error {
	if company.CompanyName == "" {
		return apperror.BadRequest("company name is required", nil)
	}
	return s.repo.Create(ctx, company)
}

func (s *Service) GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllCompanies(ctx context.Context) ([]Company, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateCompany(ctx context.Context, company *Company) error {
	if company.CompanyName == "" {
		return apperror.BadRequest("company name is required", nil)
	}
	return s.repo.Update(ctx, company)
}

func (s *Service) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
