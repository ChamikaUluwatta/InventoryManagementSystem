package validation

import (
	"strings"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ValidateProductName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return apperror.BadRequest("Invalid Product name", nil)
	}
	return nil
}

func ValidatePrice(price decimal.Decimal) error {
	if price.IsNegative() {
		return apperror.BadRequest("price cannot be negative", nil)
	}
	return nil
}

func ValidateDiameterAndWidth(diameter, width decimal.Decimal) error {
	if diameter.IsNegative() {
		return apperror.BadRequest("diameter cannot be negative", nil)
	}
	if width.IsNegative() {
		return apperror.BadRequest("width cannot be negative", nil)
	}
	return nil
}

func ValidateCompanyID(companyID uuid.UUID) error {
	if companyID == uuid.Nil {
		return apperror.BadRequest("Invalid company id", nil)
	}
	return nil
}

func ValidateProductID(productID uuid.UUID) error {
	if productID == uuid.Nil {
		return apperror.BadRequest("Invalid product id", nil)
	}
	return nil
}

func ValidateLocationID(locationID string) error {
	if locationID == "" {
		return apperror.BadRequest("Invalid location id", nil)
	}
	return nil
}

func ValidateCategoryID(categoryID int) error {
	if categoryID <= 0 {
		return apperror.BadRequest("Invalid category id", nil)
	}
	return nil
}

func ValidateCreateProduct(req *model.CreateProductRequest) error {
	if req == nil {
		return apperror.BadRequest("request body is required", nil)
	}
	if err := ValidateProductName(req.ProductName); err != nil {
		return err
	}
	if err := ValidatePrice(req.Price); err != nil {
		return err
	}
	if err := ValidateCompanyID(req.CompanyID); err != nil {
		return err
	}
	if err := ValidateCategoryID(req.CategoryID); err != nil {
		return err
	}
	if err := ValidateLocationID(req.LocationID); err != nil {
		return err
	}
	if err := ValidateDiameterAndWidth(req.Diameter, req.Width); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateProduct(req *model.Product) error {
	if req == nil {
		return apperror.BadRequest("request body is required", nil)
	}
	if err := ValidateProductID(req.ProductID); err != nil {
		return err
	}
	if err := ValidateProductName(req.ProductName); err != nil {
		return err
	}
	if err := ValidatePrice(req.Price); err != nil {
		return err
	}
	if err := ValidateCompanyID(req.CompanyID); err != nil {
		return err
	}
	if err := ValidateCategoryID(req.CategoryID); err != nil {
		return err
	}
	if err := ValidateLocationID(req.LocationID); err != nil {
		return err
	}
	if err := ValidateDiameterAndWidth(req.Diameter, req.Width); err != nil {
		return err
	}
	return nil
}
