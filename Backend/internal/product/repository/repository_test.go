package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	testDB          *testutil.TestDB
	seedCompanyID   uuid.UUID
	seedCategoryID  int
	seedLocationID  string
	seedCategoryID2 int
	seedLocationID2 string
	testTenantID    = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

const migrationPath = "../../database/migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := testutil.SetupTestDB(ctx, migrationPath)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("SetupTestDB failed: %v\n", err)
		os.Exit(1)
	}
	testDB = db

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "companies" (company_name, tenant_id) VALUES ('Test Company', $1) RETURNING company_id`,
		testTenantID,
	).Scan(&seedCompanyID); err != nil {
		fmt.Print(err.Error())
		testDB.Close()
		fmt.Printf("SetupTestDB failed: %v\n", err)
		os.Exit(1)
	}

	if err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO "categories" (category_name, tenant_id) VALUES ('Test Category', $1) RETURNING category_id`,
		testTenantID,
	).Scan(&seedCategoryID); err != nil {
		fmt.Print(err.Error())
		testDB.Close()
		fmt.Printf("SetupTestDB failed: %v\n", err)
		os.Exit(1)
	}

	seedLocationID = "TEST-LOC-1"
	if _, err := testDB.Pool.Exec(ctx,
		`INSERT INTO "locations" (location_id, tenant_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		seedLocationID, testTenantID,
	); err != nil {
		fmt.Print(err.Error())
		testDB.Close()
		fmt.Printf("SetupTestDB failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	t.Run("success", func(t *testing.T) {
		req := model.Product{
			ProductName:        "Test Create Product",
			ProductDescription: "A test product",
			Diameter:           decimal.NewFromFloat(10.0),
			Width:              decimal.NewFromFloat(5.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(9.99),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		}
		err := repo.Create(t.Context(), &req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if req.ProductID == uuid.Nil {
			t.Error("expected non-nil product ID")
		}
	})

	t.Run("foreign key violation - invalid category", func(t *testing.T) {
		req := model.Product{
			ProductName:        "Test FK Category",
			ProductDescription: "Should fail",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         99999,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		}
		err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("foreign key violation - invalid company", func(t *testing.T) {
		req := model.Product{
			ProductName:        "Test FK Company",
			ProductDescription: "Should fail",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          uuid.New(),
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		}
		err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected foreign key error, got nil")
		}
	})

	t.Run("duplicate unique constraint", func(t *testing.T) {
		req := model.Product{
			ProductName:        "Unique Product Create",
			ProductDescription: "will duplicate",
			Diameter:           decimal.NewFromFloat(3.0),
			Width:              decimal.NewFromFloat(3.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(3.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		}
		if err := repo.Create(t.Context(), &req); err != nil {
			t.Fatalf("first create should succeed, got %v", err)
		}
		err := repo.Create(t.Context(), &req)
		if err == nil {
			t.Fatal("expected unique constraint violation, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.Product{
		ProductName:        "Test GetByID Product",
		ProductDescription: "for get by id test",
		Diameter:           decimal.NewFromFloat(10.0),
		Width:              decimal.NewFromFloat(5.0),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(9.99),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
		TenantID:           testTenantID,
	}
	err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	t.Run("zero for stock when no inventory", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), req.ProductID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Stock != 0 {
			t.Errorf("expected stock 0, got %d", got.Stock)
		}
	})

	_, err = testDB.Pool.Exec(t.Context(),
		`INSERT INTO "inventories" (product_id, location_id, stock, tenant_id) VALUES ($1, $2, $3, $4)`,
		req.ProductID, seedLocationID, int32(50), testTenantID,
	)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), req.ProductID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductID != req.ProductID {
			t.Errorf("expected product ID %v, got %v", req.ProductID, got.ProductID)
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
		`INSERT INTO "categories" (category_name, tenant_id) VALUES ('Category Two', $1) RETURNING category_id`,
		testTenantID,
	).Scan(&catID2)
	if err != nil {
		t.Fatalf("failed to create second category: %v", err)
	}

	products := []model.Product{
		{
			ProductName:        "AA GetAll Product",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
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
			TenantID:           testTenantID,
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
			TenantID:           testTenantID,
		},
	}

	for i := range products {
		if err := repo.Create(t.Context(), &products[i]); err != nil {
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

	t.Run("pagination limits results", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.GetProductsQueryParams{Limit: 1, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) > 1 {
			t.Errorf("expected at most 1 product with limit 1, got %d", len(got))
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.Product{
		ProductName:        "Original Update Product",
		ProductDescription: "original description",
		Diameter:           decimal.NewFromFloat(5.0),
		Width:              decimal.NewFromFloat(2.5),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(4.99),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
		TenantID:           testTenantID,
	}
	err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	_, err = testDB.Pool.Exec(t.Context(),
		`INSERT INTO "inventories" (product_id, location_id, stock, tenant_id) VALUES ($1, $2, $3, $4)`,
		req.ProductID, seedLocationID, int32(10), testTenantID,
	)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated := model.Product{
			ProductID:          req.ProductID,
			ProductName:        "Updated Name",
			ProductDescription: "updated description",
			Diameter:           decimal.NewFromFloat(15.0),
			Width:              decimal.NewFromFloat(7.5),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(14.99),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), req.ProductID)
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
			TenantID:           testTenantID,
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
			ProductID:          req.ProductID,
			ProductName:        "Updated Product",
			ProductDescription: "Updated description",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         "non exist",
			TenantID:           testTenantID,
		}
		if err := repo.Update(t.Context(), &wrongLocation); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	req := model.Product{
		ProductName:        "Delete Me Product",
		ProductDescription: "to be deleted",
		Diameter:           decimal.NewFromFloat(1.0),
		Width:              decimal.NewFromFloat(1.0),
		CompanyID:          seedCompanyID,
		Price:              decimal.NewFromFloat(1.0),
		CategoryID:         seedCategoryID,
		LocationID:         seedLocationID,
		TenantID:           testTenantID,
	}
	err := repo.Create(t.Context(), &req)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), req.ProductID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), req.ProductID)
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

func TestProductCount(t *testing.T) {
	ctx := t.Context()
	repo := repository.NewRepository(testDB.Pool)

	// Clear products from other tests that might pollute counts
	if _, err := testDB.Pool.Exec(ctx, `DELETE FROM products`); err != nil {
		t.Fatalf("failed to clear products: %v", err)
	}

	err := testDB.Pool.QueryRow(t.Context(), `INSERT INTO categories (category_name, tenant_id) VALUES ($1, $2) RETURNING category_id`, "TestCategory2", testTenantID).Scan(&seedCategoryID2)
	if err != nil {
		t.Fatalf("failed to create categories: %v", err)
	}

	err = testDB.Pool.QueryRow(t.Context(), `INSERT INTO locations (location_id, tenant_id) VALUES ($1, $2) RETURNING location_id`, "TEST-LOC-2", testTenantID).Scan(&seedLocationID2)
	if err != nil {
		t.Fatalf("failed to create locations: %v", err)
	}

	products := []model.Product{
		{
			ProductName:        "Count Product 1",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(1.0),
			Width:              decimal.NewFromFloat(1.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(1.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		},
		{
			ProductName:        "Count Product 2",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(2.0),
			Width:              decimal.NewFromFloat(2.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(2.0),
			CategoryID:         seedCategoryID2,
			LocationID:         seedLocationID2,
			TenantID:           testTenantID,
		},
		{
			ProductName:        "Count Product 3",
			ProductDescription: "desc",
			Diameter:           decimal.NewFromFloat(3.0),
			Width:              decimal.NewFromFloat(3.0),
			CompanyID:          seedCompanyID,
			Price:              decimal.NewFromFloat(3.0),
			CategoryID:         seedCategoryID,
			LocationID:         seedLocationID,
			TenantID:           testTenantID,
		},
	}

	for i := range products {
		if err := repo.Create(t.Context(), &products[i]); err != nil {
			t.Fatalf("failed to create product %s: %v", products[i].ProductName, err)
		}
	}

	type productCount struct {
		CategoryID int `json:"category_id"`
		Count      int `json:"count"`
	}

	t.Run("Product Count Group by Category", func(t *testing.T) {
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := []productCount{
			{
				CategoryID: seedCategoryID,
				Count:      2,
			},
			{
				CategoryID: seedCategoryID2,
				Count:      1,
			},
		}

		response, err := repo.GetAllProductCountBy(t.Context(), "category")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for i, expectedGroup := range expected {
			if i >= len(response) {
				t.Fatalf("expected group %d, but response has only %d groups", i, len(response))
			}
			respGroup := response[i]
			if respGroup.CategoryID != expectedGroup.CategoryID {
				t.Errorf("expected category id %d, got %d", expectedGroup.CategoryID, respGroup.CategoryID)
			}
			if respGroup.ProductCount != expectedGroup.Count {
				t.Errorf("expected count %d for category %d, got %d", expectedGroup.Count, expectedGroup.CategoryID, respGroup.ProductCount)
			}
		}
	})

	t.Run("Product Count group by location", func(t *testing.T) {
		type locationCount struct {
			LocationID string `json:"location_id"`
			Count      int    `json:"count"`
		}
		expected := []locationCount{
			{
				LocationID: seedLocationID,
				Count:      2,
			},
			{
				LocationID: seedLocationID2,
				Count:      1,
			},
		}

		response, err := repo.GetAllProductCountBy(t.Context(), "location")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for i, expectedGroup := range expected {
			if i >= len(response) {
				t.Fatalf("expected group %d, but response has only %d groups", i, len(response))
			}
			respGroup := response[i]
			if respGroup.LocationID != expectedGroup.LocationID {
				t.Errorf("expected location id '%s', got '%s'", expectedGroup.LocationID, respGroup.LocationID)
			}
			if respGroup.ProductCount != expectedGroup.Count {
				t.Errorf("expected count %d for location '%s', got %d", expectedGroup.Count, expectedGroup.LocationID, respGroup.ProductCount)
			}
		}
	})

}
