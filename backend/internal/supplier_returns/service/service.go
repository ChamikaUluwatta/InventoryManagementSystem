package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/validation"
)

type Service interface {
	CreateSupplierReturn(ctx context.Context, req *model.SupplierReturn) error
	GetSupplierReturnByID(ctx context.Context, id int) (*model.SupplierReturn, error)
	GetAllSupplierReturns(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error)
	UpdateSupplierReturnStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error)
	DeleteSupplierReturn(ctx context.Context, id int) error
}

type service struct {
	repo repository.SupplierReturnRepository
}

func NewService(repo repository.SupplierReturnRepository) *service {
	return &service{repo: repo}
}

func (s *service) CreateSupplierReturn(ctx context.Context, req *model.SupplierReturn) error {
	if err := validation.ValidateCreateSupplierReturnRequest(req); err != nil {
		return err
	}

	return s.repo.Create(ctx, req)
}

func (s *service) GetSupplierReturnByID(ctx context.Context, id int) (*model.SupplierReturn, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllSupplierReturns(ctx context.Context, params model.QueryParams) ([]model.SupplierReturn, error) {
	validatedParams, err := validation.ValidateParams(params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, validatedParams)
}

func (s *service) UpdateSupplierReturnStatus(ctx context.Context, id int, status model.ReturnStatus) (*model.SupplierReturn, error) {
	if err := validation.ValidateUpdateSupplierReturnStatus(status); err != nil {
		return nil, err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *service) DeleteSupplierReturn(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
