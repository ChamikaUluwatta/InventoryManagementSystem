package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

func ValidateLocationID(id string) error {
	if id == "" {
		return apperror.BadRequest("location id is required", nil)
	}
	return nil
}