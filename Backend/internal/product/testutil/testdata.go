package testutil

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	CreateProductRequestMock = model.CreateProductRequest{
		ProductName:        "Test Product",
		ProductDescription: "This is a test product",
		Diameter:           decimal.NewFromFloat(10.0),
		Width:              decimal.NewFromFloat(5.0),
		CompanyID:          uuid.New(),
		Price:              decimal.NewFromFloat(9.99),
		CategoryID:         1,
		LocationID:         "A1",
	}

	ProductMock = model.Product{
		ProductID:          uuid.New(),
		ProductName:        "Test Product",
		ProductDescription: "This is a test product",
		Diameter:           decimal.NewFromFloat(10.0),
		Width:              decimal.NewFromFloat(5.0),
		CompanyID:          uuid.New(),
		Price:              decimal.NewFromFloat(9.99),
		CategoryID:         1,
		LocationID:         "A1",
	}

	GetProductByIdMock = model.GetProductById{
		Product: ProductMock,
		Stock:   100,
	}
)
