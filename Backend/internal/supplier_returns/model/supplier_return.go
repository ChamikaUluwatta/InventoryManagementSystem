package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ReturnStatus string

const (
	ReturnStatusDraft     ReturnStatus = "draft"
	ReturnStatusApproved  ReturnStatus = "approved"
	ReturnStatusSent      ReturnStatus = "sent"
	ReturnStatusCredited  ReturnStatus = "credited"
	ReturnStatusCancelled ReturnStatus = "cancelled"
	ReturnStatusRejected  ReturnStatus = "rejected"
	ReturnStatusCompleted ReturnStatus = "completed"
)

type SupplierReturn struct {
	SupplierReturnID int                  `db:"supplier_return_id" json:"supplier_return_id"`
	ReturnNo         string               `db:"return_no"          json:"return_no"`
	CompanyID        uuid.UUID            `db:"company_id"         json:"company_id"`
	Status           ReturnStatus         `db:"status"             json:"status"`
	Reason           *string              `db:"reason"             json:"reason,omitempty"`
	Notes            *string              `db:"notes"              json:"notes,omitempty"`
	CreatedAt        time.Time            `db:"created_at"         json:"created_at"`
	ApprovedAt       *time.Time           `db:"approved_at"        json:"approved_at,omitempty"`
	CompletedAt      *time.Time           `db:"completed_at"       json:"completed_at,omitempty"`
	Items            []SupplierReturnItem `db:"-" json:"items,omitempty"`
}

type SupplierReturnItem struct {
	SupplierReturnItemID int             `db:"supplier_return_item_id" json:"supplier_return_item_id"`
	SupplierReturnID     int             `db:"supplier_return_id"      json:"supplier_return_id"`
	ProductID            *uuid.UUID      `db:"product_id"              json:"product_id,omitempty"`
	LocationID           *string         `db:"location_id"             json:"location_id,omitempty"`
	Quantity             int             `db:"quantity"                json:"quantity"`
	UnitCost             decimal.Decimal `db:"unit_cost"               json:"unit_cost"`

	ProductNameSnapshot string `db:"product_name_snapshot" json:"product_name_snapshot"`
	LocationSnapshot    string `db:"location_snapshot"     json:"location_snapshot"`
}

type CreateSupplierReturnItemRequest struct {
	ProductID  *uuid.UUID      `json:"product_id,omitempty"`
	LocationID *string         `json:"location_id,omitempty"`
	Quantity   int             `json:"quantity"`
	UnitCost   decimal.Decimal `json:"unit_cost"`
}

type UpdateSupplierReturnStatusRequest struct {
	Status ReturnStatus `json:"status"`
}

func (s ReturnStatus) IsValid() bool {
	switch s {
	case ReturnStatusDraft, ReturnStatusApproved, ReturnStatusSent, ReturnStatusCredited, ReturnStatusCancelled, ReturnStatusRejected, ReturnStatusCompleted:
		return true
	default:
		return false
	}
}

type QueryParams struct {
	Limit     int
	Offset    int
	CompanyID *uuid.UUID
}
