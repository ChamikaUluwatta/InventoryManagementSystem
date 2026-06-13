package validation

import (
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
)

func ValidateLocationID(id string) error {
	if id == "" {
		return apperror.BadRequest("location id is required", nil)
	}
	return nil
}

func ValidateParams(limit, offset int) (model.QueryParams, error) {
	if limit < 0 {
		return model.QueryParams{}, apperror.BadRequest("limit must be non-negative", nil)
	}
	if offset < 0 {
		return model.QueryParams{}, apperror.BadRequest("offset must be non-negative", nil)
	}

	if limit == 0 {
		limit = 10
	}

	if limit > 100 {
		return model.QueryParams{}, apperror.BadRequest("limit must be less than or equal to 100", nil)
	}
	return model.QueryParams{Limit: limit, Offset: offset}, nil
}
