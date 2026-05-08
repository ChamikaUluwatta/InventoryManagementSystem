package validation_test

import (
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestValidateProductName(t *testing.T) {
	t.Run("Check empty name", func(t *testing.T) {
		cases := []string{
			"",
			" ",
		}
		for _, name := range cases {
			got := validation.ValidateProductName(name)
			want := errors.New("Invalid Product name")
			if got == nil {
				t.Fatalf("got %v nil want '%v'", got, want.Error())
			}
			assertEqual(t, got, want)
		}
	})
}

func TestValidatePrice(t *testing.T) {
	t.Run("Check negative price", func(t *testing.T) {
		got := validation.ValidatePrice(decimal.NewFromInt(-1))
		want := errors.New("price cannot be negative")
		if got == nil {
			t.Fatalf("got %v nil want '%v'", got, want.Error())
		}
		assertEqual(t, got, want)
	})
}

func TestDiameterAndWidthValidation(t *testing.T) {
	t.Run("Check negative diameters", func(t *testing.T) {
		got := validation.ValidateDiameterAndWidth(decimal.NewFromInt(-1), decimal.NewFromInt(10))
		want := errors.New("diameter cannot be negative")
		if got == nil {
			t.Fatalf("got %v nil want '%v'", got, want.Error())
		}
		assertEqual(t, got, want)
	})
	t.Run("Check negative widths", func(t *testing.T) {
		got := validation.ValidateDiameterAndWidth(decimal.NewFromInt(10), decimal.NewFromInt(-1))
		want := errors.New("width cannot be negative")
		if got == nil {
			t.Fatalf("got %v nil want '%v'", got, want.Error())
		}
		assertEqual(t, got, want)
	})
}

func TestCompanyIDValidation(t *testing.T) {
	t.Run("Check empty company id", func(t *testing.T) {
		got := validation.ValidateCompanyID(uuid.Nil)
		want := errors.New("Invalid company id")
		if got == nil {
			t.Fatalf("got %v nil want '%v'", got, want.Error())
		}
		assertEqual(t, got, want)
	})
}

func TestCategoryIDValidation(t *testing.T) {
	t.Run("Check empty category id", func(t *testing.T) {
		got := validation.ValidateCategoryID(0)
		want := errors.New("Invalid category id")
		if got == nil {
			t.Fatalf("got %v nil want '%v'", got, want.Error())
		}
		assertEqual(t, got, want)
	})
}

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

func assertEqual(t *testing.T, got, want error) {
	if got.Error() != want.Error() {
		t.Fatalf("got '%v' want '%v'", got.Error(), want.Error())
	}
}
