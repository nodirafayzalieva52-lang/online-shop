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
	
}

