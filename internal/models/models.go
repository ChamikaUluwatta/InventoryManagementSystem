package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID   uuid.UUID       `db:"product_id"   json:"product_id"`
	ProductName string          `db:"product_name" json:"product_name"`
	Diameter    decimal.Decimal `db:"diameter"     json:"diameter"`
	Width       decimal.Decimal `db:"width"        json:"width"`
	CompanyID   uuid.UUID       `db:"company_id"   json:"company_id"`
	Price       decimal.Decimal `db:"price"        json:"price"`
	CategoryID  int             `db:"category_id"  json:"category_id"`
}

type Category struct {
	CategoryID   int    `db:"category_id"   json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
	ParentID     *int   `db:"parent_id"     json:"parent_id"`
}

type Company struct {
	CompanyID   uuid.UUID `db:"company_id"   json:"company_id"`
	CompanyName string    `db:"company_name" json:"company_name"`
}

type Location struct {
	LocationID string  `db:"location_id" json:"location_id"`
	Image      *string `db:"image"       json:"image"`
}

type Inventory struct {
	InventoryID int       `db:"inventory_id" json:"inventory_id"`
	ProductID   uuid.UUID `db:"product_id"   json:"product_id"`
	LocationID  string    `db:"location_id"  json:"location_id"`
	Stock       int       `db:"stock"        json:"stock"`
}
