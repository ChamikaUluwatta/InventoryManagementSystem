package product

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID          uuid.UUID       `db:"product_id"   json:"product_id"`
	ProductName        string          `db:"product_name" json:"product_name"`
	ProductDescription string          `db:"product_description" json:"product_description"`
	Diameter           decimal.Decimal `db:"diameter"     json:"diameter"`
	Width              decimal.Decimal `db:"width"        json:"width"`
	CompanyID          uuid.UUID       `db:"company_id"   json:"company_id"`
	Price              decimal.Decimal `db:"price"        json:"price"`
	CategoryID         int             `db:"category_id"  json:"category_id"`
	LocationID         string          `db:"location_id"  json:"location_id"`
}

type GetProductById struct {
	Product
	Stock int `db:"stock" json:"stock"`
}
