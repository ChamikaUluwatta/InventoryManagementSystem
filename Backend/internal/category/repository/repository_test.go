package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/testutil"
)

var testDB *testutil.TestDB

const migrationPath = "../../database/migrations"

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := testutil.SetupTestDB(ctx, migrationPath)
	if err != nil {
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
		category := model.Category{
			CategoryName: "Create Test Category",
		}
		err := repo.Create(t.Context(), &category)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if category.CategoryID == 0 {
			t.Error("expected non-zero category ID")
		}
		got, err := repo.GetByID(t.Context(), category.CategoryID)
		if err != nil {
			t.Fatalf("failed to verify created category: %v", err)
		}
		if got.CategoryName != category.CategoryName {
			t.Errorf("expected name '%s', got '%s'", category.CategoryName, got.CategoryName)
		}
	})

	t.Run("success with parent", func(t *testing.T) {
		parent := model.Category{
			CategoryName: "Parent Category",
		}
		if err := repo.Create(t.Context(), &parent); err != nil {
			t.Fatalf("failed to create parent category: %v", err)
		}

		child := model.Category{
			CategoryName: "Child Category",
			ParentID:     &parent.CategoryID,
		}
		err := repo.Create(t.Context(), &child)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), child.CategoryID)
		if err != nil {
			t.Fatalf("failed to get child category: %v", err)
		}
		if got.ParentID == nil || *got.ParentID != parent.CategoryID {
			t.Errorf("expected parent ID %d, got %v", parent.CategoryID, got.ParentID)
		}
	})

	t.Run("Failure with Non exist Parent", func(t *testing.T) {
		parent := model.Category{
			CategoryName: "Real Parent Category",
		}
		if err := repo.Create(t.Context(), &parent); err != nil {
			t.Fatalf("failed to create parent category: %v", err)
		}
		nonExistentParentId := 99999
		child := model.Category{
			CategoryName: "Child Category",
			ParentID:     &nonExistentParentId,
		}
		err := repo.Create(t.Context(), &child)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "failed to create category" {
			t.Errorf("expected 'failed to create category', got '%s'", err.Error())
		}
	})
}

func TestGetByID(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	category := model.Category{
		CategoryName: "GetByID Test Category",
	}
	if err := repo.Create(t.Context(), &category); err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	t.Run("found", func(t *testing.T) {
		got, err := repo.GetByID(t.Context(), category.CategoryID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.CategoryID != category.CategoryID {
			t.Errorf("expected category ID %d, got %d", category.CategoryID, got.CategoryID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestGetAll(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	categories := []model.Category{
		{CategoryName: "BB Category"},
		{CategoryName: "AA Category"},
		{CategoryName: "CC Category"},
	}
	for i := range categories {
		if err := repo.Create(t.Context(), &categories[i]); err != nil {
			t.Fatalf("failed to create category %s: %v", categories[i].CategoryName, err)
		}
	}

	t.Run("returns all sorted by name", func(t *testing.T) {
		got, err := repo.GetAll(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var testCategories []string
		for _, c := range got {
			if c.CategoryName == "AA Category" || c.CategoryName == "BB Category" || c.CategoryName == "CC Category" {
				testCategories = append(testCategories, c.CategoryName)
			}
		}
		if len(testCategories) < 3 {
			t.Fatalf("expected at least 3 test categories, got %d", len(testCategories))
		}
		for i := 1; i < len(testCategories); i++ {
			if testCategories[i-1] > testCategories[i] {
				t.Errorf("categories not sorted by name: '%s' > '%s'", testCategories[i-1], testCategories[i])
			}
		}
	})
}

func TestUpdate(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	category := model.Category{
		CategoryName: "Original Update Category",
	}
	if err := repo.Create(t.Context(), &category); err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		updated := model.Category{
			CategoryID:   category.CategoryID,
			CategoryName: "Updated Category Name",
		}
		if err := repo.Update(t.Context(), &updated); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := repo.GetByID(t.Context(), category.CategoryID)
		if err != nil {
			t.Fatalf("failed to get updated category: %v", err)
		}
		if got.CategoryName != "Updated Category Name" {
			t.Errorf("expected 'Updated Category Name', got '%s'", got.CategoryName)
		}
	})

	t.Run("not found", func(t *testing.T) {
		nonExistent := model.Category{
			CategoryID:   99999,
			CategoryName: "Non-existent",
		}
		err := repo.Update(t.Context(), &nonExistent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestDelete(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	category := model.Category{
		CategoryName: "Delete Me Category",
	}
	if err := repo.Create(t.Context(), &category); err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := repo.Delete(t.Context(), category.CategoryID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err := repo.GetByID(t.Context(), category.CategoryID)
		if err == nil {
			t.Fatal("expected category to be deleted, but found it")
		}
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(t.Context(), 99999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestGetByParent(t *testing.T) {
	repo := repository.NewRepository(testDB.Pool)

	parent := model.Category{
		CategoryName: "Parent For GetByParent",
	}
	if err := repo.Create(t.Context(), &parent); err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	child := model.Category{
		CategoryName: "Child For GetByParent",
		ParentID:     &parent.CategoryID,
	}
	if err := repo.Create(t.Context(), &child); err != nil {
		t.Fatalf("failed to create child: %v", err)
	}

	t.Run("returns subcategories for given parent", func(t *testing.T) {
		got, err := repo.GetByParent(t.Context(), &parent.CategoryID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) == 0 {
			t.Fatal("expected subcategories, got empty")
		}
		for _, c := range got {
			if c.ParentID == nil || *c.ParentID != parent.CategoryID {
				t.Errorf("expected parent ID %d for '%s', got %v", parent.CategoryID, c.CategoryName, c.ParentID)
			}
		}
	})

	t.Run("nil parent returns root categories", func(t *testing.T) {
		got, err := repo.GetByParent(t.Context(), nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, c := range got {
			if c.ParentID != nil {
				t.Errorf("expected nil parent for '%s', got %v", c.CategoryName, c.ParentID)
			}
		}
	})

	t.Run("non-existent parent returns empty", func(t *testing.T) {
		nonExistent := 99999
		got, err := repo.GetByParent(t.Context(), &nonExistent)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %d categories", len(got))
		}
	})
}
