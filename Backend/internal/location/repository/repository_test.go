package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/location/repository"
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
		loc := model.Location{
			LocationID: "CREATE-TEST-LOC",
			TenantID:   testTenantID,
		}
		err := repo.Create(t.Context(), &loc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), loc.LocationID)
		if err != nil {
			t.Fatalf("failed to verify created location: %v", err)
		}
		if got.LocationID != loc.LocationID {
			t.Errorf("expected location ID '%s', got '%s'", loc.LocationID, got.LocationID)
		}
	})

	t.Run("success with image", func(t *testing.T) {
		img := "https://example.com/img.png"
		loc := model.Location{
			LocationID: "CREATE-WITH-IMAGE",
			Image:      &img,
			TenantID:   testTenantID,
		}
		err := repo.Create(t.Context(), &loc)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), loc.LocationID)
		if err != nil {
			t.Fatalf("failed to verify created location: %v", err)
		}
		if got.Image == nil || *got.Image != img {
			t.Errorf("expected image '%s', got '%v'", img, got.Image)
		}
	})

	t.Run("duplicate location ID returns error", func(t *testing.T) {
		loc := model.Location{
			LocationID: "DUP-LOC",
			TenantID:   testTenantID,
		}
		if err := repo.Create(t.Context(), &loc); err != nil {
			t.Fatalf("first create should succeed, got %v", err)
		}
		duplicate := model.Location{
			LocationID: "DUP-LOC",
			TenantID:   testTenantID,
		}
		err := repo.Create(t.Context(), &duplicate)
		if err == nil {
			t.Fatal("expected error for duplicate location ID, got nil")
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	loc := model.Location{
		LocationID: "GETBYID-LOC",
		TenantID:   testTenantID,
	}
	if err := repo.Create(t.Context(), &loc); err != nil {
		t.Fatalf("failed to create location: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), loc.LocationID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.LocationID != loc.LocationID {
			t.Errorf("expected location ID '%s', got '%s'", loc.LocationID, got.LocationID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), "NONEXIST-ID")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	locations := []model.Location{
		{LocationID: "B-GETALL-LOC", TenantID: testTenantID},
		{LocationID: "A-GETALL-LOC", TenantID: testTenantID},
		{LocationID: "C-GETALL-LOC", TenantID: testTenantID},
	}
	for i := range locations {
		if err := repo.Create(t.Context(), &locations[i]); err != nil {
			t.Fatalf("failed to create location %s: %v", locations[i].LocationID, err)
		}
	}

	t.Run("returns all sorted by location id", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var getAllLocs []string
		for _, l := range got {
			if l.LocationID == "A-GETALL-LOC" || l.LocationID == "B-GETALL-LOC" || l.LocationID == "C-GETALL-LOC" {
				getAllLocs = append(getAllLocs, l.LocationID)
			}
		}
		if len(getAllLocs) < 3 {
			t.Fatalf("expected at least 3 test locations, got %d", len(getAllLocs))
		}
		for i := 1; i < len(getAllLocs); i++ {
			if getAllLocs[i-1] > getAllLocs[i] {
				t.Errorf("locations not sorted by ID: '%s' > '%s'", getAllLocs[i-1], getAllLocs[i])
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	loc := model.Location{
		LocationID: "UPDATE-LOC",
		TenantID:   testTenantID,
	}
	if err := repo.Create(t.Context(), &loc); err != nil {
		t.Fatalf("failed to create location: %v", err)
	}

	t.Run("success with image", func(t *testing.T) {
		img := "https://example.com/updated.png"
		updated := model.Location{
			LocationID: loc.LocationID,
			Image:      &img,
			TenantID:   testTenantID,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), loc.LocationID)
		if err != nil {
			t.Fatalf("failed to get updated location: %v", err)
		}
		if got.Image == nil || *got.Image != img {
			t.Errorf("expected image '%s', got '%v'", img, got.Image)
		}
	})

	t.Run("success with nil image", func(t *testing.T) {
		updated := model.Location{
			LocationID: loc.LocationID,
			Image:      nil,
			TenantID:   testTenantID,
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), loc.LocationID)
		if err != nil {
			t.Fatalf("failed to get updated location: %v", err)
		}
		if got.Image != nil {
			t.Errorf("expected nil image, got '%s'", *got.Image)
		}
	})

	t.Run("not found", func(t *testing.T) {
		nonExistent := model.Location{
			LocationID: "NONEXIST-UPDATE",
			TenantID:   testTenantID,
		}
		err := repo.Update(t.Context(), &nonExistent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	loc := model.Location{
		LocationID: "DELETE-LOC",
		TenantID:   testTenantID,
	}
	if err := repo.Create(t.Context(), &loc); err != nil {
		t.Fatalf("failed to create location: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), loc.LocationID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), loc.LocationID)
		if err == nil {
			t.Fatal("expected location to be deleted, but found it")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), "NONEXIST-DELETE")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "location not found" {
			t.Errorf("expected 'location not found', got '%s'", err.Error())
		}
	})
}
