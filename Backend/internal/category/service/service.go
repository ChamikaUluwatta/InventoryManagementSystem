package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/validation"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, category *model.Category) error {
	if err := validation.ValidateCategoryName(category.CategoryName); err != nil {
		return err
	}
	return s.repo.Create(ctx, category)
}

func (s *Service) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := validation.ValidateCategoryName(category.CategoryName); err != nil {
		return err
	}
	return s.repo.Update(ctx, category)
}

func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCategoriesByParent(ctx context.Context, parentID *int) ([]model.Category, error) {
	return s.repo.GetByParent(ctx, parentID)
}