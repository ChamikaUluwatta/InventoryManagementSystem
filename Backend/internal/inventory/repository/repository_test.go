package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

var (
	testDB         *testutil.TestDB
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

	var companyID uuid.UUID
	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "companies" (company_name) VALUES ('Test Company') RETURNING company_id`,
	).Scan(&companyID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	var categoryID int
	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "categories" (category_name) VALUES ('Test Category') RETURNING category_id`,
	).Scan(&categoryID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	seedLocationID = "INV-LOC-1"
	if _, err := testDB.Pool.Exec(ctx,
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		seedLocationID,
	); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "products" (product_name, product_description, diameter, width, company_id, price, category_id, location_id)
		 VALUES ('Test Product', 'desc', 1.0, 1.0, $1, 1.0, $2, $3) RETURNING product_id`,
		companyID, categoryID, seedLocationID,
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

	t.Run("success", func(t *testing.T) {
		inv := model.Inventory{
			ProductID:  seedProductID,
			LocationID: seedLocationID,
			Stock:      50,
		}
		err := repo.Create(t.Context(), &inv)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if inv.InventoryID == 0 {
			t.Error("expected non-zero inventory ID")
		}
	})

	t.Run("foreign key violation - invalid product", func(t *testing.T) {
		inv := model.Inventory{
			ProductID:  uuid.New(),
			LocationID: seedLocationID,
			Stock:      10,
		}
		_, err := repo.GetByID(t.Context(), 99999)
		if err == nil {
			t.Log("skipping FK test - need valid approach")
			return
		}

		err = repo.Create(t.Context(), &inv)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("foreign key violation - invalid location", func(t *testing.T) {
		inv := model.Inventory{
			ProductID:  seedProductID,
			LocationID: "NONEXIST-LOC",
			Stock:      10,
		}
		err := repo.Create(t.Context(), &inv)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("unique constraint violation", func(t *testing.T) {
		inv := model.Inventory{
			ProductID:  seedProductID,
			LocationID: seedLocationID,
			Stock:      25,
		}
		err := repo.Create(t.Context(), &inv)
		if err == nil {
			t.Fatal("expected unique constraint violation, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	inv := model.Inventory{
		ProductID:  seedProductID,
		LocationID: "GETBYID-LOC",
		Stock:      75,
	}
	if _, err := testDB.Pool.Exec(t.Context(),
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		"GETBYID-LOC",
	); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}
	if err := repo.Create(t.Context(), &inv); err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), inv.InventoryID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.InventoryID != inv.InventoryID {
			t.Errorf("expected inventory ID %d, got %d", inv.InventoryID, got.InventoryID)
		}
		if got.Stock != 75 {
			t.Errorf("expected stock 75, got %d", got.Stock)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	t.Run("returns inventories sorted by id", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(got); i++ {
			if got[i-1].InventoryID > got[i].InventoryID {
				t.Errorf("inventories not sorted by ID: %d > %d", got[i-1].InventoryID, got[i].InventoryID)
			}
		}
	})

	t.Run("filters by product", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{
			Limit:     10,
			Offset:    0,
			ProductID: &seedProductID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, inv := range got {
			if inv.ProductID != seedProductID {
				t.Errorf("expected product_id %v, got %v", seedProductID, inv.ProductID)
			}
		}
	})

	t.Run("filters by location", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{
			Limit:      10,
			Offset:     0,
			LocationID: &seedLocationID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, inv := range got {
			if inv.LocationID != seedLocationID {
				t.Errorf("expected location_id %s, got %s", seedLocationID, inv.LocationID)
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	locID := "UPDATE-LOC"
	if _, err := testDB.Pool.Exec(t.Context(),
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		locID,
	); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}

	inv := model.Inventory{
		ProductID:  seedProductID,
		LocationID: locID,
		Stock:      30,
	}
	if err := repo.Create(t.Context(), &inv); err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated := model.Inventory{
			InventoryID: inv.InventoryID,
			ProductID:   seedProductID,
			LocationID:  locID,
			Stock:       500,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), inv.InventoryID)
		if err != nil {
			t.Fatalf("failed to get updated inventory: %v", err)
		}
		if got.Stock != 500 {
			t.Errorf("expected stock 500, got %d", got.Stock)
		}
	})

	t.Run("not found", func(t *testing.T) {
		nonExistent := model.Inventory{
			InventoryID: 99999,
			ProductID:   seedProductID,
			LocationID:  locID,
			Stock:       10,
		}
		err := repo.Update(t.Context(), &nonExistent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	locID := "DELETE-LOC"
	if _, err := testDB.Pool.Exec(t.Context(),
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		locID,
	); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}

	inv := model.Inventory{
		ProductID:  seedProductID,
		LocationID: locID,
		Stock:      10,
	}
	if err := repo.Create(t.Context(), &inv); err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), inv.InventoryID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), inv.InventoryID)
		if err == nil {
			t.Fatal("expected inventory to be deleted, but found it")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "inventory not found" {
			t.Errorf("expected 'inventory not found', got '%s'", err.Error())
		}
	})
}


