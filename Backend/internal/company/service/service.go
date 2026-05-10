package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/validation"
	"github.com/google/uuid"
)

type Service interface {
	CreateCompany(ctx context.Context, company *model.Company) error
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error)
	GetAllCompanies(ctx context.Context, params model.QueryParams) ([]model.Company, error)
	UpdateCompany(ctx context.Context, company *model.Company) error
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	GetCompanyDependencies(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateCompany(ctx context.Context, company *model.Company) error {
	if err := validation.ValidateCompanyName(company.CompanyName); err != nil {
		return err
	}
	return s.repo.Create(ctx, company)
}

func (s *service) GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllCompanies(ctx context.Context, params model.QueryParams) ([]model.Company, error) {
	validatedParams, err := validation.ValidateParams(params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, validatedParams)
}

func (s *service) UpdateCompany(ctx context.Context, company *model.Company) error {
	if err := validation.ValidateCompanyName(company.CompanyName); err != nil {
		return err
	}
	return s.repo.Update(ctx, company)
}

func (s *service) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetCompanyDependencies(ctx context.Context, id uuid.UUID) (model.CompanyDependency, error) {
	return s.repo.GetCompanyDependencies(ctx, id)
}
