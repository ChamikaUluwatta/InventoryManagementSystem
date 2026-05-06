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

func ValidateParams(limit, offset int) (error, model.QueryParams) {
	if limit < 0 {
		return apperror.BadRequest("limit must be non-negative", nil), model.QueryParams{}
	}
	if offset < 0 {
		return apperror.BadRequest("offset must be non-negative", nil), model.QueryParams{}
	}

	if limit == 0 {
		limit = 10
	}

	if limit > 100 {
		return apperror.BadRequest("limit must be less than or equal to 100", nil), model.QueryParams{}
	}
	return nil, model.QueryParams{Limit: limit, Offset: offset}
}
