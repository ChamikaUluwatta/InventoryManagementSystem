package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/shopspring/decimal"
)

func ValidateProductName(name string) error {
	if name == "" {
		return apperror.BadRequest("product name is required", nil)
	}
	return nil
}

func ValidatePrice(price decimal.Decimal) error {
	if price.IsNegative() {
		return apperror.BadRequest("price cannot be negative", nil)
	}
	return nil
}