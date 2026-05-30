package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

var testDB *testutil.TestDB
var testTenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

const migrationPath = "../../database/migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := testutil.SetupTestDB(ctx, migrationPath)
	if err != nil {
		fmt.Printf("SetupTestDB failed: %v\n", err)
		os.Exit(1)
	}
	testDB = db

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	t.Run("success", func(t *testing.T) {
		company := model.Company{
			CompanyName: "Create Test Company",
			TenantID:    testTenantID,
		}
		err := repo.Create(t.Context(), &company)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if company.CompanyID == uuid.Nil {
			t.Error("expected non-nil company ID")
		}

		got, err := repo.GetByID(t.Context(), company.CompanyID)
		if err != nil {
			t.Fatalf("failed to verify created company: %v", err)
		}
		if got.CompanyName != company.CompanyName {
			t.Errorf("expected name '%s', got '%s'", company.CompanyName, got.CompanyName)
		}
	})

	t.Run("duplicate name returns error", func(t *testing.T) {
		company := model.Company{
			CompanyName: "Duplicate Name Company",
			TenantID:    testTenantID,
		}
		if err := repo.Create(t.Context(), &company); err != nil {
			t.Fatalf("first create should succeed, got %v", err)
		}
		duplicate := model.Company{
			CompanyName: "Duplicate Name Company",
			TenantID:    testTenantID,
		}
		err := repo.Create(t.Context(), &duplicate)
		if err == nil {
			t.Fatal("expected error for duplicate company name, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	company := model.Company{
		CompanyName: "GetByID Test Company",
		TenantID:    testTenantID,
	}
	if err := repo.Create(t.Context(), &company); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), company.CompanyID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.CompanyID != company.CompanyID {
			t.Errorf("expected company ID %v, got %v", company.CompanyID, got.CompanyID)
		}
		if got.CompanyName != company.CompanyName {
			t.Errorf("expected name '%s', got '%s'", company.CompanyName, got.CompanyName)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	companies := []model.Company{
		{CompanyName: "BB Test Company", TenantID: testTenantID},
		{CompanyName: "AA Test Company", TenantID: testTenantID},
		{CompanyName: "CC Test Company", TenantID: testTenantID},
	}
	for i := range companies {
		if err := repo.Create(t.Context(), &companies[i]); err != nil {
			t.Fatalf("failed to create company %s: %v", companies[i].CompanyName, err)
		}
	}

	t.Run("returns all sorted by name", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var testCompanies []string
		for _, c := range got {
			if c.CompanyName == "AA Test Company" || c.CompanyName == "BB Test Company" || c.CompanyName == "CC Test Company" {
				testCompanies = append(testCompanies, c.CompanyName)
			}
		}
		if len(testCompanies) < 3 {
			t.Fatalf("expected at least 3 test companies, got %d", len(testCompanies))
		}
		for i := 1; i < len(testCompanies); i++ {
			if testCompanies[i-1] > testCompanies[i] {
				t.Errorf("companies not sorted by name: '%s' > '%s'", testCompanies[i-1], testCompanies[i])
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	company := model.Company{
		CompanyName: "Original Update Company",
		TenantID:    testTenantID,
	}
	if err := repo.Create(t.Context(), &company); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated := model.Company{
			CompanyID:   company.CompanyID,
			CompanyName: "Updated Company Name",
			TenantID:    testTenantID,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), company.CompanyID)
		if err != nil {
			t.Fatalf("failed to get updated company: %v", err)
		}
		if got.CompanyName != "Updated Company Name" {
			t.Errorf("expected 'Updated Company Name', got '%s'", got.CompanyName)
		}
	})

	t.Run("not found", func(t *testing.T) {
		nonExistent := model.Company{
			CompanyID:   uuid.New(),
			CompanyName: "Non-existent",
			TenantID:    testTenantID,
		}
		err := repo.Update(t.Context(), &nonExistent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	company := model.Company{
		CompanyName: "Delete Me Company",
		TenantID:    testTenantID,
	}
	if err := repo.Create(t.Context(), &company); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), company.CompanyID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), company.CompanyID)
		if err == nil {
			t.Fatal("expected company to be deleted, but found it")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "company not found" {
			t.Errorf("expected 'company not found', got '%s'", err.Error())
		}
	})
}

func TestGetCompanyDependencies(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	company := model.Company{
		CompanyName: "Test Company",
		TenantID:    testTenantID,
	}
	if err := repo.Create(t.Context(), &company); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	companyZeroDeps := model.Company{
		CompanyName: "Zero Deps Company",
		TenantID:    testTenantID,
	}
	if err := repo.Create(t.Context(), &companyZeroDeps); err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	var categoryID int
	if err := testDB.Pool.QueryRow(t.Context(), `INSERT INTO "categories" (category_name, tenant_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING category_id`, "Test Category", testTenantID).Scan(&categoryID); err != nil {
		t.Fatalf("failed to create category dependency: %v", err)
	}

	if _, err := testDB.Pool.Exec(t.Context(), `INSERT INTO "products" (product_name,company_id,category_id,tenant_id) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`, "Test Product", company.CompanyID, categoryID, testTenantID); err != nil {
		t.Fatalf("failed to create product dependency: %v", err)
	}

	if _, err := testDB.Pool.Exec(t.Context(), `INSERT INTO "supplier_returns" (return_no,company_id,tenant_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`, "Test Return", company.CompanyID, testTenantID); err != nil {
		t.Fatalf("failed to create supplier return dependency: %v", err)
	}

	t.Run("Get dependencies when exist", func(t *testing.T) {
		got, err := repo.GetCompanyDependencies(t.Context(), company.CompanyID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductCount != 1 {
			t.Errorf("expected 1 product, got %d", got.ProductCount)
		}
		if got.SupplierCount != 1 {
			t.Errorf("expected 1 supplier, got %d", got.SupplierCount)
		}
	})

	t.Run("Get zero count when No Dependencies", func(t *testing.T) {
		got, err := repo.GetCompanyDependencies(t.Context(), companyZeroDeps.CompanyID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ProductCount != 0 {
			t.Errorf("expected 0 products, got %d", got.ProductCount)
		}
		if got.SupplierCount != 0 {
			t.Errorf("expected 0 suppliers, got %d", got.SupplierCount)
		}
	})

}
