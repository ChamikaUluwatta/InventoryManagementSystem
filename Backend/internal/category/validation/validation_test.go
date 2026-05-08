package validation_test

import (
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/validation"
)

func TestValidateParams(t *testing.T) {
	t.Run("negative limit returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(-1, 0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be non-negative" {
			t.Errorf("expected 'limit must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("negative offset returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(10, -1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "offset must be non-negative" {
			t.Errorf("expected 'offset must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("zero limit defaults to 10", func(t *testing.T) {
		params, err := validation.ValidateParams(0, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if params.Limit != 10 {
			t.Errorf("expected default limit 10, got %d", params.Limit)
		}
	})
	t.Run("limit exceeds max returns error", func(t *testing.T) {
		_, err := validation.ValidateParams(101, 0)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})
	t.Run("valid limit and offset returns params", func(t *testing.T) {
		params, err := validation.ValidateParams(25, 50)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if params.Limit != 25 {
			t.Errorf("expected limit 25, got %d", params.Limit)
		}
		if params.Offset != 50 {
			t.Errorf("expected offset 50, got %d", params.Offset)
		}
	})
}

func TestValidateCategoryName(t *testing.T) {
	t.Run("empty name returns error", func(t *testing.T) {
		got := validation.ValidateCategoryName("")
		want := errors.New("category name is required")
		if got == nil {
			t.Fatalf("expected error, got nil")
		}
		if got.Error() != want.Error() {
			t.Fatalf("got '%v' want '%v'", got.Error(), want.Error())
		}
	})

	t.Run("non-empty name returns nil", func(t *testing.T) {
		got := validation.ValidateCategoryName("Electronics")
		if got != nil {
			t.Fatalf("expected nil, got '%v'", got.Error())
		}
	})
}
