package service

import (
	"context"
	"fmt"

	"shop/internal/models"
	"shop/internal/repository"
)

type StoreService struct {
	storeRepo repository.StoreRepository
}

func NewStoreService(storeRepo repository.StoreRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
	}
}

func (s *StoreService) CreateStore(ctx context.Context, sellerID int, name, description string) (*models.Store, error) {
	if name == "" {
		return nil, fmt.Errorf("название магазина не может быть пустым")
	}

	store := &models.Store{
		Seller_ID: int64(sellerID),
		Name:        name,
		Description: description,
	}

	err := s.storeRepo.Create(ctx, store)
	if err != nil {
		return nil, fmt.Errorf("storeRepo.Create: %w", err)
	}

	return store, nil
}

func (s *StoreService) GetByID(ctx context.Context, id int) (*models.Store, error) {
	if id <= 0 {
		return nil, fmt.Errorf("некорректный id магазина")
	}

	store, err := s.storeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("магазин не найден: %w", err)
	}

	return store, nil
}

func (s *StoreService) GetBySellerID(ctx context.Context, sellerID int) (*models.Store, error) {
	if sellerID <= 0 {
		return nil, fmt.Errorf("некорректный id продавца")
	}

	store, err := s.storeRepo.GetBySellerID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("магазин продавца не найден: %w", err)
	}

	return store, nil
}