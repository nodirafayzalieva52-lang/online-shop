package service

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"
	"shop/internal/repository"
	appErr "shop/pkg/errors"
)

type OrderService struct {
	OrderRepo repository.OrderRepository
	StoreRepo repository.StoreRepository
}

func NewOrderService(orderRepo repository.OrderRepository, storeRepo repository.StoreRepository) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
		StoreRepo: storeRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order) error {
	if len(order.Items) == 0 {
		return appErr.ErrEmptyOrder
	}
	if order.StoreID <= 0 {
		return errors.New("store_id is required")
	}

	if err := s.OrderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (s *OrderService) GetByID(ctx context.Context, orderID int64, userID int64, userRole domain.Role) (*domain.Order, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order id")
	}

	order, err := s.OrderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if order == nil {
		return nil, appErr.ErrOrderNotFound
	}

	// Permission check: customer who placed order, seller who owns store, or admin
	if userRole != domain.RoleAdmin && order.CustomerID != userID {
		store, err := s.StoreRepo.GetByID(ctx, order.StoreID)
		if err != nil || store == nil || store.SellerID != userID {
			return nil, appErr.ErrAccessDenied
		}
	}

	return order, nil
}

func (s *OrderService) GetByCustomerID(ctx context.Context, customerID int64) ([]*domain.Order, error) {
	if customerID <= 0 {
		return nil, errors.New("invalid customer id")
	}

	return s.OrderRepo.GetByCustomerID(ctx, customerID)
}

func (s *OrderService) GetByStoreID(ctx context.Context, storeID int64, userID int64, userRole domain.Role) ([]*domain.Order, error) {
	if storeID <= 0 {
		return nil, errors.New("invalid store id")
	}

	if userRole != domain.RoleAdmin {
		store, err := s.StoreRepo.GetByID(ctx, storeID)
		if err != nil || store == nil {
			return nil, appErr.ErrStoreNotFound
		}
		if store.SellerID != userID {
			return nil, appErr.ErrAccessDenied
		}
	}

	return s.OrderRepo.GetByStoreID(ctx, storeID)
}

func (s *OrderService) UpdateStatus(ctx context.Context, orderID int64, userID int64, userRole domain.Role, newStatus domain.OrderStatus) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	order, err := s.OrderRepo.GetByID(ctx, orderID)
	if err != nil || order == nil {
		return appErr.ErrOrderNotFound
	}

	if userRole != domain.RoleAdmin {
		store, err := s.StoreRepo.GetByID(ctx, order.StoreID)
		if err != nil || store == nil || store.SellerID != userID {
			return appErr.ErrAccessDenied
		}
	}

	return s.OrderRepo.UpdateStatus(ctx, orderID, newStatus)
}
