package service

import (
	"context"
	"fmt"

	"shop/internal/models"
	"shop/internal/repository"
)

type ProductService struct {
	ProductRepo repository.ProductRespository
}

func NewProductService(ProductRepo repository.ProductRespository) *ProductService {
	return &ProductService{
		ProductRepo: ProductRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, storeID int, name, description string, price float64, stock int) (*models.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("название товара не может быть пустым")
	}
	if price <= 0 {
		return nil, fmt.Errorf("цена товара должна быть больше нуля")
	}
	if stock < 0 {
		return nil, fmt.Errorf("количество товара не может быть отрицательным")
	}

	product := &models.Product{
		StoreID:     storeID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	err := s.ProductRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("productRepo.Create: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("некорректный id товара")
	}

	product, err := s.ProductRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("товар не найден: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetByStoreID(ctx context.Context, storeID int) ([]*models.Product, error) {
	if storeID <= 0 {
		return nil, fmt.Errorf("некорректный id магазина")
	}

	products, err := s.ProductRepo.GetByStoreID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении товаров магазина: %w", err)
	}

	return products, nil
}
