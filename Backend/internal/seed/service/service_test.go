package service_test

import (
	"context"
	"os"
	"testing"

	categoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	companyRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/company/repository"
	inventoryRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/inventory/repository"
	locationRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
	productRepo "github.com/ChamikaUluwatta/Inventory_Management_System/internal/product/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/seed/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
	"github.com/google/uuid"
)

var testDB *testutil.TestDB
var testTenantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

const migrationsDir = "../../database/migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, err := testutil.SetupTestDB(ctx, migrationsDir)
	if err != nil {
		os.Exit(1)
	}
	testDB = db
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func newSeedService() *service.Service {
	return service.NewService(
		companyRepo.NewRepository(testDB.Pool),
		categoryRepo.NewRepository(testDB.Pool),
		locationRepo.NewRepository(testDB.Pool),
		productRepo.NewRepository(testDB.Pool),
		inventoryRepo.NewRepository(testDB.Pool),
		testDB.Pool,
	)
}

func TestSeed_CreatesAllEntitiesWithCorrectCounts(t *testing.T) {
	svc := newSeedService()

	result, ids, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("Seed() returned unexpected error: %v", err)
	}

	if result.Companies != 3 {
		t.Errorf("expected companies_created = 3, got %d", result.Companies)
	}
	if result.Categories != 3 {
		t.Errorf("expected categories_created = 3, got %d", result.Categories)
	}
	if result.Locations != 4 {
		t.Errorf("expected locations_created = 4, got %d", result.Locations)
	}
	if result.Products != 3 {
		t.Errorf("expected products_created = 3, got %d", result.Products)
	}
	if result.Inventories != 3 {
		t.Errorf("expected inventories_created = 3, got %d", result.Inventories)
	}

	if len(ids.CompanyIDs) != 3 {
		t.Errorf("expected 3 company IDs, got %d", len(ids.CompanyIDs))
	}
	if len(ids.CategoryIDs) != 3 {
		t.Errorf("expected 3 category IDs, got %d", len(ids.CategoryIDs))
	}
	if len(ids.LocationIDs) != 4 {
		t.Errorf("expected 4 location IDs, got %d", len(ids.LocationIDs))
	}
	if len(ids.ProductIDs) != 3 {
		t.Errorf("expected 3 product IDs, got %d", len(ids.ProductIDs))
	}
}

func TestSeed_PersistsDataToDatabase(t *testing.T) {
	svc := newSeedService()

	_, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("Seed() failed: %v", err)
	}

	var companyCount int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM companies").Scan(&companyCount)
	if companyCount != 3 {
		t.Errorf("expected 3 rows in companies, got %d", companyCount)
	}

	var categoryCount int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM categories").Scan(&categoryCount)
	if categoryCount != 3 {
		t.Errorf("expected 3 rows in categories, got %d", categoryCount)
	}

	var locationCount int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM locations").Scan(&locationCount)
	if locationCount != 4 {
		t.Errorf("expected 4 rows in locations, got %d", locationCount)
	}

	var productCount int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM products").Scan(&productCount)
	if productCount != 3 {
		t.Errorf("expected 3 rows in products, got %d", productCount)
	}

	var inventoryCount int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM inventories").Scan(&inventoryCount)
	if inventoryCount != 3 {
		t.Errorf("expected 3 rows in inventories, got %d", inventoryCount)
	}
}

func TestSeed_ForeignKeyReferencesAreValid(t *testing.T) {
	svc := newSeedService()

	_, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("Seed() failed: %v", err)
	}

	var productRefCount int
	err = testDB.Pool.QueryRow(t.Context(), `
		SELECT COUNT(*) FROM products p
		JOIN companies c ON c.company_id = p.company_id
		JOIN categories cat ON cat.category_id = p.category_id
	`).Scan(&productRefCount)
	if err != nil {
		t.Fatalf("failed to query product FK refs: %v", err)
	}
	if productRefCount != 3 {
		t.Errorf("expected 3 products with valid company/category FK refs, got %d", productRefCount)
	}

	var invRefCount int
	err = testDB.Pool.QueryRow(t.Context(), `
		SELECT COUNT(*) FROM inventories i
		JOIN products p ON p.product_id = i.product_id
		JOIN locations l ON l.location_id = i.location_id
	`).Scan(&invRefCount)
	if err != nil {
		t.Fatalf("failed to query inventory FK refs: %v", err)
	}
	if invRefCount != 3 {
		t.Errorf("expected 3 inventories with valid product/location FK refs, got %d", invRefCount)
	}
}

func TestSeed_RerunClearsAndReSeeds(t *testing.T) {
	svc := newSeedService()

	_, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("first Seed() failed: %v", err)
	}

	_, err = testDB.Pool.Exec(t.Context(), "INSERT INTO companies (company_name, tenant_id) VALUES ('Foreign Corp', $1)", testTenantID)
	if err != nil {
		t.Fatalf("failed to insert extra row: %v", err)
	}

	var countBefore int
	testDB.Pool.QueryRow(t.Context(), "SELECT COUNT(*) FROM companies").Scan(&countBefore)
	if countBefore != 4 {
		t.Fatalf("expected 4 companies before re-seed, got %d", countBefore)
	}

	result, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("second Seed() failed: %v", err)
	}

	if result.Companies != 3 {
		t.Errorf("expected 3 companies after re-seed, got %d", result.Companies)
	}
	if result.Products != 3 {
		t.Errorf("expected 3 products after re-seed, got %d", result.Products)
	}
	if result.Inventories != 3 {
		t.Errorf("expected 3 inventories after re-seed, got %d", result.Inventories)
	}
}

func TestSeed_TruncateAndReSeed(t *testing.T) {
	svc := newSeedService()

	_, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("first Seed() failed: %v", err)
	}

	_, err = testDB.Pool.Exec(t.Context(), "TRUNCATE TABLE inventories, products, locations, categories, companies RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	result, _, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("Seed() after truncate failed: %v", err)
	}

	if result.Companies != 3 {
		t.Errorf("expected 3 companies after re-seed, got %d", result.Companies)
	}
	if result.Products != 3 {
		t.Errorf("expected 3 products after re-seed, got %d", result.Products)
	}
	if result.Inventories != 3 {
		t.Errorf("expected 3 inventories after re-seed, got %d", result.Inventories)
	}
}

func TestSeed_ReturnsNonEmptyIDs(t *testing.T) {
	svc := newSeedService()

	_, ids, err := svc.Seed(t.Context(), testTenantID)
	if err != nil {
		t.Fatalf("Seed() failed: %v", err)
	}

	for i, id := range ids.CompanyIDs {
		if id == uuid.Nil {
			t.Errorf("company_ids[%d] is nil UUID", i)
		}
	}
	for i, id := range ids.CategoryIDs {
		if id == 0 {
			t.Errorf("category_ids[%d] is 0", i)
		}
	}
	for i, id := range ids.LocationIDs {
		if id == "" {
			t.Errorf("location_ids[%d] is empty", i)
		}
	}
	for i, id := range ids.ProductIDs {
		if id == uuid.Nil {
			t.Errorf("product_ids[%d] is nil UUID", i)
		}
	}
}
