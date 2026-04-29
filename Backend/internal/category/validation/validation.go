package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func ValidateCategoryName(name string) error {
	if name == "" {
		return apperror.BadRequest("category name is required", nil)
	}
	return nil
}