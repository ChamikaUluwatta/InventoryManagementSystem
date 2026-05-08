package testutil

import (
	"time"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func SupplierReturnMock() model.SupplierReturn {
	return model.SupplierReturn{
		SupplierReturnID: 1,
		ReturnNo:         "SR-001",
		CompanyID:        uuid.New(),
		Status:           model.ReturnStatusDraft,
		CreatedAt:        time.Now(),
	}
}

func SupplierReturnWithItemsMock() model.SupplierReturn {
	return model.SupplierReturn{
		SupplierReturnID: 1,
		ReturnNo:         "SR-002",
		CompanyID:        uuid.New(),
		Status:           model.ReturnStatusDraft,
		CreatedAt:        time.Now(),
		Items: []model.SupplierReturnItem{
			SupplierReturnItemMock(),
		},
	}
}

func SupplierReturnItemMock() model.SupplierReturnItem {
	productID := uuid.New()
	locationID := "TEST-LOC-1"
	return model.SupplierReturnItem{
		SupplierReturnItemID: 1,
		SupplierReturnID:     1,
		ProductID:            &productID,
		LocationID:           &locationID,
		Quantity:             5,
		UnitCost:             decimal.NewFromFloat(10.50),
		ProductNameSnapshot:  "Test Product",
		LocationSnapshot:     "TEST-LOC-1",
	}
}

func CreateSupplierReturnMock() model.SupplierReturn {
	productID := uuid.New()
	locationID := "TEST-LOC-1"
	return model.SupplierReturn{
		ReturnNo:  "SR-NEW",
		CompanyID: uuid.New(),
		Reason:    strPtr("Defective items"),
		Notes:     strPtr("Urgent return"),
		Items: []model.SupplierReturnItem{
			{
				ProductID:  &productID,
				LocationID: &locationID,
				Quantity:   10,
				UnitCost:   decimal.NewFromFloat(5.99),
			},
		},
	}
}

func UpdateSupplierReturnStatusRequestMock() model.UpdateSupplierReturnStatusRequest {
	return model.UpdateSupplierReturnStatusRequest{
		Status: model.ReturnStatusApproved,
	}
}

func strPtr(s string) *string {
	return &s
}
