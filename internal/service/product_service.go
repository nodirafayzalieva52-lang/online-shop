package service

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/delivery/http/dto"
	"shop/internal/domain"
	"shop/internal/repository"
	appErr "shop/pkg/errors"
)

type ProductService struct {
	productRepo repository.ProductRepository
	storeRepo   repository.StoreRepository
}

func NewProductService(productRepo repository.ProductRepository, storeRepo repository.StoreRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, userID int64, userRole domain.Role, req dto.CreateProductRequest) (*domain.Product, error) {
	if req.Name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if req.Price <= 0 {
		return nil, errors.New("product price must be greater than zero")
	}
	if req.Stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	store, err := s.storeRepo.GetByID(ctx, req.StoreID)
	if err != nil || store == nil {
		return nil, appErr.ErrStoreNotFound
	}

	if store.SellerID != userID && userRole != domain.RoleAdmin {
		return nil, appErr.ErrAccessDenied
	}

	product := &domain.Product{
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	err = s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("productRepo.Create: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product id")
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return nil, appErr.ErrProductNotFound
	}

	return product, nil
}

func (s *ProductService) GetByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]*domain.Product, error) {
	if storeID <= 0 {
		return nil, errors.New("invalid store id")
	}

	products, err := s.productRepo.GetByStoreID(ctx, storeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store products: %w", err)
	}

	return products, nil
}

func (s *ProductService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	return s.productRepo.GetAll(ctx, limit, offset)
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID int64, userID int64, userRole domain.Role, req dto.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	store, err := s.storeRepo.GetByID(ctx, product.StoreID)
	if err != nil || store == nil {
		return nil, appErr.ErrStoreNotFound
	}

	if store.SellerID != userID && userRole != domain.RoleAdmin {
		return nil, appErr.ErrAccessDenied
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != nil && *req.Name != "" {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("product price must be greater than zero")
		}
		product.Price = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("product stock cannot be negative")
		}
		product.Stock = *req.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64, userID int64, userRole domain.Role) error {
	product, err := s.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	store, err := s.storeRepo.GetByID(ctx, product.StoreID)
	if err != nil || store == nil {
		return appErr.ErrStoreNotFound
	}

	if store.SellerID != userID && userRole != domain.RoleAdmin {
		return appErr.ErrAccessDenied
	}

	return s.productRepo.Delete(ctx, productID)
}
