package service

import (
	"context"

	"shop/internal/models"
	"shop/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{Name: name}
	err := s.categoryRepo.Create(ctx, *category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*models.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*models.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}