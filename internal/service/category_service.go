package service

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"
	"shop/internal/repository"
	appErr "shop/pkg/errors"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category := &domain.Category{Name: name}
	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("categoryRepo.Create: %w", err)
	}
	return category, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid category id")
	}
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, appErr.ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*domain.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int64, name string) (*domain.Category, error) {
	category, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category.Name = name
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("categoryRepo.Update: %w", err)
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, id)
}
