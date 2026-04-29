package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func ValidateStock(stock int) error {
	if stock < 0 {
		return apperror.BadRequest("stock cannot be negative", nil)
	}
	return nil
}