package service

import (
	"context"

	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/model"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/repository"
	"github.com/ChamikaUluwatta/Inventory_Management_System/internal/category/validation"
)

type Service interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, id int) (*model.Category, error)
	GetAllCategories(ctx context.Context, params model.QueryParams) ([]model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategoriesByParent(ctx context.Context, parentID *int) ([]model.Category, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(ctx context.Context, category *model.Category) error {
	if err := validation.ValidateCategoryName(category.CategoryName); err != nil {
		return err
	}
	return s.repo.Create(ctx, category)
}

func (s *service) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAllCategories(ctx context.Context, params model.QueryParams) ([]model.Category, error) {
	validatedParams, err := validation.ValidateParams(params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, validatedParams)
}

func (s *service) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := validation.ValidateCategoryName(category.CategoryName); err != nil {
		return err
	}
	return s.repo.Update(ctx, category)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetCategoriesByParent(ctx context.Context, parentID *int) ([]model.Category, error) {
	return s.repo.GetByParent(ctx, parentID)
}