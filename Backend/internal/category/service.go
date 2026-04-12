package category

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/apperror"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, category *Category) error {
	if category.CategoryName == "" {
		return apperror.BadRequest("category name is required", nil)
	}
	return s.repo.Create(ctx, category)
}

func (s *Service) GetCategoryByID(ctx context.Context, id int) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateCategory(ctx context.Context, category *Category) error {
	if category.CategoryName == "" {
		return apperror.BadRequest("category name is required", nil)
	}
	return s.repo.Update(ctx, category)
}

func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCategoriesByParent(ctx context.Context, parentID *int) ([]Category, error) {
	return s.repo.GetByParent(ctx, parentID)
}
