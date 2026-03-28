package product

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
