package validation

import (
	"strings"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/google/uuid"
)

func ValidateCreateSupplierReturnRequest(req *model.CreateSupplierReturnRequest) error {
	if req == nil {
		return apperror.BadRequest("request body is required", nil)
	}

	req.ReturnNo = strings.TrimSpace(req.ReturnNo)
	if req.ReturnNo == "" {
		return apperror.BadRequest("return_no is required", nil)
	}

	if req.CompanyID == uuid.Nil {
		return apperror.BadRequest("company_id is required", nil)
	}

	if len(req.Items) == 0 {
		return apperror.BadRequest("at least one item is required", nil)
	}

	for i := range req.Items {
		item := &req.Items[i]

		if item.ProductID == nil {
			return apperror.BadRequest("item product_id is required", nil)
		}
		if *item.ProductID == uuid.Nil {
			return apperror.BadRequest("item product_id is invalid", nil)
		}
		if item.LocationID == nil {
			return apperror.BadRequest("item location_id is required", nil)
		}
		*item.LocationID = strings.TrimSpace(*item.LocationID)
		if *item.LocationID == "" {
			return apperror.BadRequest("item location_id is required", nil)
		}

		if item.Quantity <= 0 {
			return apperror.BadRequest("item quantity must be greater than 0", nil)
		}
		if item.UnitCost.IsNegative() {
			return apperror.BadRequest("item unit_cost cannot be negative", nil)
		}
	}

	return nil
}

func ValidateUpdateSupplierReturnStatus(status model.ReturnStatus) error {
	if !status.IsValid() {
		return apperror.BadRequest("invalid supplier return status", nil)
	}
	return nil
}