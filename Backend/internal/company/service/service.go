package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCompany(ctx context.Context, company *model.Company) error {
	if err := validation.ValidateCompanyName(company.CompanyName); err != nil {
		return err
	}
	return s.repo.Create(ctx, company)
}

func (s *Service) GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllCompanies(ctx context.Context) ([]model.Company, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateCompany(ctx context.Context, company *model.Company) error {
	if err := validation.ValidateCompanyName(company.CompanyName); err != nil {
		return err
	}
	return s.repo.Update(ctx, company)
}

func (s *Service) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}