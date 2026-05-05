package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	testDB         *testutil.TestDB
	seedCompanyID  uuid.UUID
	seedCategoryID int
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
		`INSERT INTO "companies" (company_name) VALUES ('Test Company') RETURNING company_id`,
	).Scan(&seedCompanyID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "categories" (category_name) VALUES ('Test Category') RETURNING category_id`,
	).Scan(&seedCategoryID); err != nil {
		testDB.Close()
		os.Exit(1)
	}

	seedLocationID = "TEST-LOC-1"
	if _, err := testDB.Pool.Exec(ctx,
		`INSERT INTO "locations" (location_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		seedLocationID,
	); err != nil {
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
		req := model.CreateProductRequest{
			ProductName:        "Test Create Product",
			ProductDescription: "A test product",
			Diameter:           decimal.NewFromFloat(10.0),
			Width:              decimal.NewFromFloat(5.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(9.99),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		}
		product, err := repo.Create(t.Context(), &req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if product.ProductID == uuid.Nil {
			t.Error("expected non-nil product ID")
		}
	})

	t.Run("foreign key violation - invalid category", func(t *testing.T) {
		req := model.CreateProductRequest{
			ProductName:        "Test FK Category",
			ProductDescription: "Should fail",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         99999,
			LocationID:         seedLocationID,
		}
		_, err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("foreign key violation - invalid company", func(t *testing.T) {
		req := model.CreateProductRequest{
			ProductName:        "Test FK Company",
			ProductDescription: "Should fail",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          uuid.New(),
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		}
		_, err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("duplicate unique constraint", func(t *testing.T) {
		req := model.CreateProductRequest{
			ProductName:        "Unique Product Create",
			ProductDescription: "will duplicate",
			Diameter:           decimal.NewFromFloat(3.0),
			Width:              decimal.NewFromFloat(3.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(3.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		}
		if _, err := repo.Create(t.Context(), &req); err != nil {
			t.Fatalf("first create should succeed, got %v", err)
		}
		_, err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected unique constraint violation, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.CreateProductRequest{
		ProductName:        "Test GetByID Product",
		ProductDescription: "for get by id test",
		Diameter:           decimal.NewFromFloat(10.0),
		Width:              decimal.NewFromFloat(5.0),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(9.99),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
	}
	created, err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	t.Run("zero for stock when no inventory", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), created.ProductID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Stock != 0 {
			t.Errorf("expected stock 0, got %d", got.Stock)
		}
	})

	_, err = testDB.Pool.Exec(t.Context(),
		`INSERT INTO "inventories" (product_id, location_id, stock) VALUES ($1, $2, $3)`,
		created.ProductID, seedLocationID, int32(50),
	)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), created.ProductID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductID != created.ProductID {
			t.Errorf("expected product ID %v, got %v", created.ProductID, got.ProductID)
		}
		if got.Stock != 50 {
			t.Errorf("expected stock 50, got %d", got.Stock)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("expected 'product not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	var catID2 int
	err := testDB.Pool.QueryRow(t.Context(),
		`INSERT INTO "categories" (category_name) VALUES ('Category Two') RETURNING category_id`,
	).Scan(&catID2)
	if err != nil {
		t.Fatalf("failed to create second category: %v", err)
	}

	products := []model.CreateProductRequest{
		{
			ProductName:        "AA GetAll Product",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		},
		{
			ProductName:        "BB GetAll Product",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(2.0),
			Width:              decimal.NewFromFloat(2.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(2.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		},
		{
			ProductName:        "CC Different Category",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(3.0),
			Width:              decimal.NewFromFloat(3.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(3.0),
			CategoryID:         catID2,
			LocationID:         seedLocationID,
		},
	}

	for i := range products {
		if _, err := repo.Create(t.Context(), &products[i]); err != nil {
			t.Fatalf("failed to create product %s: %v", products[i].ProductName, err)
		}
	}

	t.Run("no filters", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.GetProductsQueryParams{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(got); i++ {
			if got[i-1].ProductName > got[i].ProductName {
				t.Errorf("products not sorted by name: %q > %q",
					got[i-1].ProductName, got[i].ProductName)
			}
		}
	})

	t.Run("filter by category", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.GetProductsQueryParams{
			CategoryID: &seedCategoryID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, p := range got {
			if p.CategoryID != seedCategoryID {
				t.Errorf("expected category %d, got %d for product %s",
					seedCategoryID, p.CategoryID, p.ProductName)
			}
		}
	})

	t.Run("filter by company", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.GetProductsQueryParams{
			CompanyID: &seedCompanyID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, p := range got {
			if p.CompanyID != seedCompanyID {
				t.Errorf("expected company %v, got %v for product %s",
					seedCompanyID, p.CompanyID, p.ProductName)
			}
		}
	})

	t.Run("empty result", func(t *testing.T) {
		unknownID := 99999
		got, err := repo.GetAll(t.Context(), model.GetProductsQueryParams{
			CategoryID: &unknownID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %d products", len(got))
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.CreateProductRequest{
		ProductName:        "Original Update Product",
		ProductDescription: "original description",
		Diameter:           decimal.NewFromFloat(5.0),
		Width:              decimal.NewFromFloat(2.5),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(4.99),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
	}
	created, err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	_, err = testDB.Pool.Exec(t.Context(),
		`INSERT INTO "inventories" (product_id, location_id, stock) VALUES ($1, $2, $3)`,
		created.ProductID, seedLocationID, int32(10),
	)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated := model.Product{
			ProductID:          created.ProductID,
			ProductName:        "Updated Name",
			ProductDescription: "updated description",
			Diameter:           decimal.NewFromFloat(15.0),
			Width:              decimal.NewFromFloat(7.5),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(14.99),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), created.ProductID)
		if err != nil {
			t.Fatalf("failed to get updated product: %v", err)
		}
		if got.ProductName != "Updated Name" {
			t.Errorf("expected 'Updated Name', got '%s'", got.ProductName)
		}
		if got.ProductDescription != "updated description" {
			t.Errorf("expected 'updated description', got '%s'", got.ProductDescription)
		}
		if !got.Diameter.Equal(decimal.NewFromFloat(15.0)) {
			t.Errorf("expected diameter 15.0, got %s", got.Diameter)
		}
		if !got.Price.Equal(decimal.NewFromFloat(14.99)) {
			t.Errorf("expected price 14.99, got %s", got.Price)
		}
		if got.LocationID != seedLocationID {
			t.Errorf("expected location '%s', got '%s'", seedLocationID, got.LocationID)
		}
	})

	t.Run("not found when product id is invalid", func(t *testing.T) {
		nonExistent := model.Product{
			ProductID:          uuid.New(),
			ProductName:        "Non-existent",
			ProductDescription: "no",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
		}
		err := repo.Update(t.Context(), &nonExistent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("expected 'product not found', got '%s'", err.Error())
		}
	})

	t.Run("Non exist location id", func(t *testing.T) {
		wrongLocation := model.Product{
			ProductID:          created.ProductID,
			ProductName:        "Updated Product",
			ProductDescription: "Updated description",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         "non exist",
		}
		if err := repo.Update(t.Context(), &wrongLocation); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.CreateProductRequest{
		ProductName:        "Delete Me Product",
		ProductDescription: "to be deleted",
		Diameter:           decimal.NewFromFloat(1.0),
		Width:              decimal.NewFromFloat(1.0),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(1.0),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
	}
	created, err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), created.ProductID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), created.ProductID)
		if err == nil {
			t.Fatal("expected product to be deleted, but found it")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("expected 'product not found', got '%s'", err.Error())
		}
	})
}
