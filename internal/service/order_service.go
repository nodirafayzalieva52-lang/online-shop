package service

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/models"
	"shop/internal/repository"
	"shop/pkg/errors"
)

type OrderService struct {
	OrderRepo repository.OrderRepository
}

func NewOrderService(OrderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: OrderRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if len(order.Items) == 0 {
		return pkg.ErrEmptyOrder
	}

	var totalPrice float64
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return fmt.Errorf("%w: item quantity must be greater than zero", pkg.ErrInvalidOrder)
		}
		if item.Price < 0 {
			return fmt.Errorf("%w: item price cannot be negative", pkg.ErrInvalidOrder)
		}
		totalPrice += item.Price * float64(item.Quantity)
	}

	order.TotalPrice = totalPrice
	order.Status = "created"

	if err := s.OrderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("failed to create order in repo: %w", err)
	}

	return nil
}

func (s *OrderService) GetByCustomerID(ctx context.Context, customerID int) ([]*models.Order, error) {
	if customerID <= 0 {
		return nil, errors.New("invalid customer id")
	}

	return s.OrderRepo.GetByCustomerID(ctx, customerID)
}

func (s *OrderService) GetByStoreID(ctx context.Context, storeID int) ([]*models.Order, error) {
	if storeID <= 0 {
		return nil, errors.New("invalid store id")
	}

	return s.OrderRepo.GetByStoreID(ctx, storeID)
}

func (s *OrderService) GetByID(ctx context.Context, id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order id")
	}

	return s.OrderRepo.GetByID(ctx, id)
}