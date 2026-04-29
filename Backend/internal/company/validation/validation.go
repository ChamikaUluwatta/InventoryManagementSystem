package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func ValidateCompanyName(name string) error {
	if name == "" {
		return apperror.BadRequest("company name is required", nil)
	}
	return nil
}