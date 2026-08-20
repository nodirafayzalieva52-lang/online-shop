package service

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"
	"shop/internal/repository"
	appErr "shop/pkg/errors"
)

type StoreService struct {
	storeRepo repository.StoreRepository
	userRepo  repository.UserRepository
}

func NewStoreService(storeRepo repository.StoreRepository, userRepo repository.UserRepository) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
		userRepo:  userRepo,
	}
}

func (s *StoreService) CreateStore(ctx context.Context, sellerID int64, name, description string) (*domain.Store, error) {
	if name == "" {
		return nil, errors.New("store name cannot be empty")
	}

	user, err := s.userRepo.GetByID(ctx, sellerID)
	if err != nil || user == nil {
		return nil, appErr.ErrUserNotFound
	}

	// Auto-promote customer to seller upon creating first store
	if user.Role == domain.RoleCustomer {
		if err := s.userRepo.UpdateRole(ctx, sellerID, domain.RoleSeller); err != nil {
			return nil, fmt.Errorf("failed to promote user to seller: %w", err)
		}
	}

	store := &domain.Store{
		SellerID:    sellerID,
		Name:        name,
		Description: description,
	}

	err = s.storeRepo.Create(ctx, store)
	if err != nil {
		return nil, fmt.Errorf("storeRepo.Create: %w", err)
	}

	return store, nil
}

func (s *StoreService) GetByID(ctx context.Context, id int64) (*domain.Store, error) {
	if id <= 0 {
		return nil, errors.New("invalid store id")
	}

	store, err := s.storeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store: %w", err)
	}
	if store == nil {
		return nil, appErr.ErrStoreNotFound
	}

	return store, nil
}

func (s *StoreService) GetBySellerID(ctx context.Context, sellerID int64) (*domain.Store, error) {
	if sellerID <= 0 {
		return nil, errors.New("invalid seller id")
	}

	store, err := s.storeRepo.GetBySellerID(ctx, sellerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch seller store: %w", err)
	}
	if store == nil {
		return nil, appErr.ErrStoreNotFound
	}

	return store, nil
}

func (s *StoreService) UpdateStore(ctx context.Context, storeID int64, userID int64, userRole domain.Role, name, description *string) (*domain.Store, error) {
	store, err := s.GetByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	if store.SellerID != userID && userRole != domain.RoleAdmin {
		return nil, appErr.ErrAccessDenied
	}

	if name != nil && *name != "" {
		store.Name = *name
	}
	if description != nil {
		store.Description = *description
	}

	if err := s.storeRepo.Update(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	return store, nil
}

func (s *StoreService) DeleteStore(ctx context.Context, storeID int64, userID int64, userRole domain.Role) error {
	store, err := s.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	if store.SellerID != userID && userRole != domain.RoleAdmin {
		return appErr.ErrAccessDenied
	}

	return s.storeRepo.Delete(ctx, storeID)
}
