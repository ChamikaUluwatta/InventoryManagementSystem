package model

import "github.com/google/uuid"

type Company struct {
	CompanyID   uuid.UUID `db:"company_id"   json:"company_id"`
	CompanyName string    `db:"company_name" json:"company_name"`
	Description string    `db:"description"  json:"description"`
}

type QueryParams struct {
	Limit  int
	Offset int
}

type CompanyDependency struct {
	ProductCount  int `json:"product_count"`
	SupplierCount int `json:"supplier_count"`
}
