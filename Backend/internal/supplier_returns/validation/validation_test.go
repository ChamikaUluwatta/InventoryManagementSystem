package validation_test

import (
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestValidateCreateSupplierReturnRequest(t *testing.T) {
	t.Run("nil request returns error", func(t *testing.T) {
		err := validation.ValidateCreateSupplierReturnRequest(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "request body is required" {
			t.Errorf("expected 'request body is required', got '%s'", err.Error())
		}
	})

	t.Run("empty return_no returns error", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "",
			CompanyID: uuid.New(),
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "return_no is required" {
			t.Errorf("expected 'return_no is required', got '%s'", err.Error())
		}
	})

	t.Run("whitespace return_no returns error", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "   ",
			CompanyID: uuid.New(),
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "return_no is required" {
			t.Errorf("expected 'return_no is required', got '%s'", err.Error())
		}
	})

	t.Run("nil company_id returns error", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.Nil,
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company_id is required" {
			t.Errorf("expected 'company_id is required', got '%s'", err.Error())
		}
	})

	t.Run("empty items returns error", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items:     []model.SupplierReturnItem{},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "at least one item is required" {
			t.Errorf("expected 'at least one item is required', got '%s'", err.Error())
		}
	})

	t.Run("nil item product_id returns error", func(t *testing.T) {
		locationID := "LOC-1"
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  nil,
					LocationID: &locationID,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item product_id is required" {
			t.Errorf("expected 'item product_id is required', got '%s'", err.Error())
		}
	})

	t.Run("nil uuid item product_id returns error", func(t *testing.T) {
		locationID := "LOC-1"
		nilUUID := uuid.Nil
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &nilUUID,
					LocationID: &locationID,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item product_id is invalid" {
			t.Errorf("expected 'item product_id is invalid', got '%s'", err.Error())
		}
	})

	t.Run("nil item location_id returns error", func(t *testing.T) {
		productID := uuid.New()
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: nil,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item location_id is required" {
			t.Errorf("expected 'item location_id is required', got '%s'", err.Error())
		}
	})

	t.Run("empty item location_id returns error", func(t *testing.T) {
		productID := uuid.New()
		emptyLoc := ""
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: &emptyLoc,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item location_id is required" {
			t.Errorf("expected 'item location_id is required', got '%s'", err.Error())
		}
	})

	t.Run("zero quantity returns error", func(t *testing.T) {
		productID := uuid.New()
		loc := "LOC-1"
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: &loc,
					Quantity:   0,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item quantity must be greater than 0" {
			t.Errorf("expected 'item quantity must be greater than 0', got '%s'", err.Error())
		}
	})

	t.Run("negative quantity returns error", func(t *testing.T) {
		productID := uuid.New()
		loc := "LOC-1"
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: &loc,
					Quantity:   -5,
					UnitCost:   decimal.NewFromFloat(1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item quantity must be greater than 0" {
			t.Errorf("expected 'item quantity must be greater than 0', got '%s'", err.Error())
		}
	})

	t.Run("negative unit_cost returns error", func(t *testing.T) {
		productID := uuid.New()
		loc := "LOC-1"
		req := &model.SupplierReturn{
			ReturnNo:  "SR-001",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: &loc,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(-1),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "item unit_cost cannot be negative" {
			t.Errorf("expected 'item unit_cost cannot be negative', got '%s'", err.Error())
		}
	})

	t.Run("valid request returns nil", func(t *testing.T) {
		productID := uuid.New()
		loc := "LOC-1"
		req := &model.SupplierReturn{
			ReturnNo:  "SR-VALID",
			CompanyID: uuid.New(),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &productID,
					LocationID: &loc,
					Quantity:   1,
					UnitCost:   decimal.NewFromFloat(10),
				},
			},
		}
		err := validation.ValidateCreateSupplierReturnRequest(req)
		if err != nil {
			t.Fatalf("expected nil, got '%v'", err.Error())
		}
	})
}

func TestValidateUpdateSupplierReturnStatus(t *testing.T) {
	t.Run("valid status returns nil", func(t *testing.T) {
		for _, status := range []model.ReturnStatus{
			model.ReturnStatusDraft,
			model.ReturnStatusApproved,
			model.ReturnStatusSent,
			model.ReturnStatusCredited,
			model.ReturnStatusCancelled,
			model.ReturnStatusRejected,
			model.ReturnStatusCompleted,
		} {
			err := validation.ValidateUpdateSupplierReturnStatus(status)
			if err != nil {
				t.Errorf("expected nil for status '%s', got '%v'", status, err.Error())
			}
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		err := validation.ValidateUpdateSupplierReturnStatus("invalid_status")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "invalid supplier return status" {
			t.Errorf("expected 'invalid supplier return status', got '%s'", err.Error())
		}
	})
}

func TestValidateParams(t *testing.T) {
	t.Run("negative limit returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(-1, 0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be non-negative" {
			t.Errorf("expected 'limit must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("negative offset returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(10, -1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "offset must be non-negative" {
			t.Errorf("expected 'offset must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("zero limit defaults to 10", func(t *testing.T) {
		params, err := validation.ValidateParams(0, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if params.Limit != 10 {
			t.Errorf("expected default limit 10, got %d", params.Limit)
		}
	})
	t.Run("limit exceeds max returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(101, 0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})
	t.Run("valid limit and offset returns params", func(t *testing.T) {
		params, err := validation.ValidateParams(25, 50)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if params.Limit != 25 {
			t.Errorf("expected limit 25, got %d", params.Limit)
		}
		if params.Offset != 50 {
			t.Errorf("expected offset 50, got %d", params.Offset)
		}
	})
}
