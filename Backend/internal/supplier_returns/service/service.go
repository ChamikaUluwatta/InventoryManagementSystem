package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/validation"
)

type Service struct {
	repo repository.SupplierReturnRepository
}

func NewService(repo repository.SupplierReturnRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSupplierReturn(ctx context.Context, req *model.CreateSupplierReturnRequest) (*model.SupplierReturn, error) {
	if err := validation.ValidateCreateSupplierReturnRequest(req); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) GetSupplierReturnByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllSupplierReturns(ctx context.Context) ([]model.SupplierReturn, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateSupplierReturnStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
	if err := validation.ValidateUpdateSupplierReturnStatus(status); err != nil {
		return nil, err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) DeleteSupplierReturn(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}