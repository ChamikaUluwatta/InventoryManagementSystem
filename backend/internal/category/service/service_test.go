package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/service"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/testutil"
)

type mockRepo struct {
	createFunc      func(ctx context.Context, category *model.Category) error
	getByIDFunc     func(ctx context.Context, id int) (*model.Category, error)
	getAllFunc      func(ctx context.Context, params model.QueryParams) ([]model.Category, error)
	updateFunc      func(ctx context.Context, category *model.Category) error
	deleteFunc      func(ctx context.Context, id int) error
	getByParentFunc func(ctx context.Context, parentID *int) ([]model.Category, error)
}

func (m *mockRepo) Create(ctx context.Context, category *model.Category) error {
	return m.createFunc(ctx, category)
}
func (m *mockRepo) GetByID(ctx context.Context, id int) (*model.Category, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepo) GetAll(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
	return m.getAllFunc(ctx, params)
}
func (m *mockRepo) Update(ctx context.Context, category *model.Category) error {
	return m.updateFunc(ctx, category)
}
func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}
func (m *mockRepo) GetByParent(ctx context.Context, parentID *int) ([]model.Category, error) {
	return m.getByParentFunc(ctx, parentID)
}

const emptyNameError = "category name is required"

func TestCreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		category := testutil.CategoryMock()
		err := svc.CreateCategory(t.Context(), &category)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("success with parent", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		category := testutil.CategoryWithParentMock()
		err := svc.CreateCategory(t.Context(), &category)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty name returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			createFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		category := model.Category{CategoryName: ""}
		err := svc.CreateCategory(t.Context(), &category)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != emptyNameError {
			t.Errorf("expected '%s', got '%s'", emptyNameError, err.Error())
		}
	})
}

func TestGetCategoryByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := testutil.CategoryMock()
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.Category, error) {
				return &expected, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCategoryByID(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.CategoryID != expected.CategoryID {
			t.Errorf("expected category ID %d, got %d", expected.CategoryID, got.CategoryID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			getByIDFunc: func(ctx context.Context, id int) (*model.Category, error) {
				return nil, errors.New("category not found")
			},
		}
		svc := service.NewService(mock)
		_, err := svc.GetCategoryByID(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestGetAllCategories(t *testing.T) {
	t.Run("success with results", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				return []model.Category{testutil.CategoryMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := 	svc.GetAllCategories(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 category, got %d", len(got))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := &mockRepo{
			getAllFunc: func(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
				return []model.Category{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetAllCategories(t.Context(), model.QueryParams{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 categories, got %d", len(got))
		}
	})
	t.Run("negative limit returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCategories(t.Context(), model.QueryParams{Limit: -1, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be non-negative" {
			t.Errorf("expected 'limit must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("negative offset returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCategories(t.Context(), model.QueryParams{Limit: 10, Offset: -1})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "offset must be non-negative" {
			t.Errorf("expected 'offset must be non-negative', got '%s'", err.Error())
		}
	})
	t.Run("limit over max returns validation error", func(t *testing.T) {
		mock := &mockRepo{}
		svc := service.NewService(mock)
		_, err := svc.GetAllCategories(t.Context(), model.QueryParams{Limit: 101, Offset: 0})
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != "limit must be less than or equal to 100" {
			t.Errorf("expected 'limit must be less than or equal to 100', got '%s'", err.Error())
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		category := testutil.CategoryMock()
		err := svc.UpdateCategory(t.Context(), &category)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("empty name returns validation error", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Category) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		category := model.Category{CategoryID: 1, CategoryName: ""}
		err := svc.UpdateCategory(t.Context(), &category)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
		if err.Error() != emptyNameError {
			t.Errorf("expected '%s', got '%s'", emptyNameError, err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			updateFunc: func(ctx context.Context, c *model.Category) error {
				return errors.New("category not found")
			},
		}
		svc := service.NewService(mock)
		category := testutil.CategoryMock()
		err := svc.UpdateCategory(t.Context(), &category)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestDeleteCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteCategory(t.Context(), 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockRepo{
			deleteFunc: func(ctx context.Context, id int) error {
				return errors.New("category not found")
			},
		}
		svc := service.NewService(mock)
		err := svc.DeleteCategory(t.Context(), 999)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("expected 'category not found', got '%s'", err.Error())
		}
	})
}

func TestGetCategoriesByParent(t *testing.T) {
	parentID := 1

	t.Run("success with subcategories", func(t *testing.T) {
		mock := &mockRepo{
			getByParentFunc: func(ctx context.Context, parentID *int) ([]model.Category, error) {
				return []model.Category{testutil.CategoryWithParentMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCategoriesByParent(t.Context(), &parentID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 subcategory, got %d", len(got))
		}
	})

	t.Run("nil parent returns root categories", func(t *testing.T) {
		mock := &mockRepo{
			getByParentFunc: func(ctx context.Context, parentID *int) ([]model.Category, error) {
				return []model.Category{testutil.RootCategoryMock()}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCategoriesByParent(t.Context(), nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1 root category, got %d", len(got))
		}
	})

	t.Run("no subcategories returns empty list", func(t *testing.T) {
		mock := &mockRepo{
			getByParentFunc: func(ctx context.Context, parentID *int) ([]model.Category, error) {
				return []model.Category{}, nil
			},
		}
		svc := service.NewService(mock)
		got, err := svc.GetCategoriesByParent(t.Context(), &parentID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 subcategories, got %d", len(got))
		}
	})
}
