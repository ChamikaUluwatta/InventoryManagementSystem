package supplierreturns

import (
	"context"
	"strings"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/google/uuid"
)

type Service struct {
	repo SupplierReturnRepository
}

func NewService(repo SupplierReturnRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSupplierReturn(ctx context.Context, req *CreateSupplierReturnRequest) (*SupplierReturn, error) {
	if req == nil {
		return nil, apperror.BadRequest("request body is required", nil)
	}

	req.ReturnNo = strings.TrimSpace(req.ReturnNo)
	if req.ReturnNo == "" {
		return nil, apperror.BadRequest("return_no is required", nil)
	}

	if req.CompanyID == uuid.Nil {
		return nil, apperror.BadRequest("company_id is required", nil)
	}

	if len(req.Items) == 0 {
		return nil, apperror.BadRequest("at least one item is required", nil)
	}

	for i := range req.Items {
		item := &req.Items[i]

		if item.ProductID == nil {
			return nil, apperror.BadRequest("item product_id is required", nil)
		}
		if *item.ProductID == uuid.Nil {
			return nil, apperror.BadRequest("item product_id is invalid", nil)
		}
		if item.LocationID == nil {
			return nil, apperror.BadRequest("item location_id is required", nil)
		}
		*item.LocationID = strings.TrimSpace(*item.LocationID)
		if *item.LocationID == "" {
			return nil, apperror.BadRequest("item location_id is required", nil)
		}

		if item.Quantity <= 0 {
			return nil, apperror.BadRequest("item quantity must be greater than 0", nil)
		}
		if item.UnitCost.IsNegative() {
			return nil, apperror.BadRequest("item unit_cost cannot be negative", nil)
		}
	}

	return s.repo.Create(ctx, req)
}

func (s *Service) GetSupplierReturnByID(ctx context.Context, id int) (*SupplierReturn, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllSupplierReturns(ctx context.Context) ([]SupplierReturn, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateSupplierReturnStatus(ctx context.Context, id int, status ReturnStatus) (*SupplierReturn, error) {
	if !status.IsValid() {
		return nil, apperror.BadRequest("invalid supplier return status", nil)
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) DeleteSupplierReturn(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
