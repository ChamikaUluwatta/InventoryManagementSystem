package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/supplier_returns/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	testDB         *testutil.TestDB
	seedCompanyID  uuid.UUID
	seedProductID  uuid.UUID
	seedLocationID string
)

const migrationPath = "../../database/migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := testutil.SetupTestDB(ctx, migrationPath)
	if err != nil {
		os.Exit(1)
	}
	testDB = db

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "companies" (company_name) VALUES ('Supplier Co') RETURNING company_id`,
	).Scan(&seedCompanyID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	seedLocationID = "SR-LOC-1"
	if _, err := testDB.Pool.Exec(ctx,
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		seedLocationID,
	); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	var categoryID int
	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "categories" (category_name) VALUES ('SR Category') RETURNING category_id`,
	).Scan(&categoryID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "products" (product_name, product_description, diameter, width, company_id, price, category_id, location_id)
		 VALUES ('Supplier Product', 'desc', 1.0, 1.0, $1, 10.0, $2, $3) RETURNING product_id`,
		seedCompanyID, categoryID, seedLocationID,
	).Scan(&seedProductID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	t.Run("success with items", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "SR-CREATE-001",
			CompanyID: seedCompanyID,
			Reason:    strPtr("Defective products"),
			Notes:     strPtr("Handle ASAP"),
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &seedProductID,
					LocationID: &seedLocationID,
					Quantity:   5,
					UnitCost:   decimal.NewFromFloat(9.99),
				},
			},
		}
		err := repo.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if req.SupplierReturnID == 0 {
			t.Error("expected non-zero supplier return ID")
		}

		got, err := repo.GetByID(t.Context(), req.SupplierReturnID)
		if err != nil {
			t.Fatalf("failed to get created return: %v", err)
		}
		if got.ReturnNo != "SR-CREATE-001" {
			t.Errorf("expected return_no 'SR-CREATE-001', got '%s'", got.ReturnNo)
		}
		if got.Status != model.ReturnStatusDraft {
			t.Errorf("expected status 'draft', got '%s'", got.Status)
		}
		if len(got.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(got.Items))
		}
	})

	t.Run("created return has snapshot data", func(t *testing.T) {
		req := &model.SupplierReturn{
			ReturnNo:  "SR-SNAPSHOT-001",
			CompanyID: seedCompanyID,
			Items: []model.SupplierReturnItem{
				{
					ProductID:  &seedProductID,
					LocationID: &seedLocationID,
					Quantity:   3,
					UnitCost:   decimal.NewFromFloat(4.50),
				},
			},
		}
		err := repo.Create(t.Context(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), req.SupplierReturnID)
		if err != nil {
			t.Fatalf("failed to get created return: %v", err)
		}
		if len(got.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(got.Items))
		}
		if got.Items[0].ProductNameSnapshot != "Supplier Product" {
			t.Errorf("expected snapshot 'Supplier Product', got '%s'", got.Items[0].ProductNameSnapshot)
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := &model.SupplierReturn{
		ReturnNo:  "SR-GETBYID-001",
		CompanyID: seedCompanyID,
		Items: []model.SupplierReturnItem{
			{
				ProductID:  &seedProductID,
				LocationID: &seedLocationID,
				Quantity:   7,
				UnitCost:   decimal.NewFromFloat(3.50),
			},
		},
	}
	err := repo.Create(t.Context(), req)
	if err != nil {
		t.Fatalf("failed to create return: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), req.SupplierReturnID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ReturnNo != "SR-GETBYID-001" {
			t.Errorf("expected return_no 'SR-GETBYID-001', got '%s'", got.ReturnNo)
		}
		if got.Status != model.ReturnStatusDraft {
			t.Errorf("expected status 'draft', got '%s'", got.Status)
		}
		if len(got.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(got.Items))
		}
		if got.Items[0].Quantity != 7 {
			t.Errorf("expected item quantity 7, got %d", got.Items[0].Quantity)
		}
		if got.Items[0].ProductNameSnapshot != "Supplier Product" {
			t.Errorf("expected product_name_snapshot 'Supplier Product', got '%s'", got.Items[0].ProductNameSnapshot)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "supplier return not found" {
			t.Errorf("expected 'supplier return not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	t.Run("returns returns sorted by created_at desc", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(got); i++ {
			prev := got[i-1]
			curr := got[i]
			if prev.CreatedAt.Before(curr.CreatedAt) {
				t.Errorf("returns not sorted by created_at DESC")
			}
		}
	})
}

func TestUpdateStatus(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := &model.SupplierReturn{
		ReturnNo:  "SR-STATUS-001",
		CompanyID: seedCompanyID,
		Items: []model.SupplierReturnItem{
			{
				ProductID:  &seedProductID,
				LocationID: &seedLocationID,
				Quantity:   1,
				UnitCost:   decimal.NewFromFloat(1.0),
			},
		},
	}
	err := repo.Create(t.Context(), req)
	if err != nil {
		t.Fatalf("failed to create return: %v", err)
	}

	t.Run("draft to approved sets approved_at", func(t *testing.T) {
		result, err := repo.UpdateStatus(t.Context(), req.SupplierReturnID, model.ReturnStatusApproved)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != model.ReturnStatusApproved {
			t.Errorf("expected status 'approved', got '%s'", result.Status)
		}
		if result.ApprovedAt == nil {
			t.Error("expected approved_at to be set")
		}
	})

	t.Run("approved to completed sets completed_at", func(t *testing.T) {
		result, err := repo.UpdateStatus(t.Context(), req.SupplierReturnID, model.ReturnStatusCompleted)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Status != model.ReturnStatusCompleted {
			t.Errorf("expected status 'completed', got '%s'", result.Status)
		}
		if result.CompletedAt == nil {
			t.Error("expected completed_at to be set")
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := &model.SupplierReturn{
		ReturnNo:  "SR-DELETE-001",
		CompanyID: seedCompanyID,
		Items: []model.SupplierReturnItem{
			{
				ProductID:  &seedProductID,
				LocationID: &seedLocationID,
				Quantity:   1,
				UnitCost:   decimal.NewFromFloat(1.0),
			},
		},
	}
	err := repo.Create(t.Context(), req)
	if err != nil {
		t.Fatalf("failed to create return: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), req.SupplierReturnID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "supplier return not found" {
			t.Errorf("expected 'supplier return not found', got '%s'", err.Error())
		}
	})
}

func strPtr(s string) *string {
	return &s
}
