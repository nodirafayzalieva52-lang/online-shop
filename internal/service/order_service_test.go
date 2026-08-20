package service_test

import (
	"context"
	"testing"

	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if o := args.Get(0); o != nil {
		return o.(*domain.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderRepo) GetByCustomerID(ctx context.Context, customerID int64) ([]*domain.Order, error) {
	args := m.Called(ctx, customerID)
	if o := args.Get(0); o != nil {
		return o.([]*domain.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderRepo) GetByStoreID(ctx context.Context, storeID int64) ([]*domain.Order, error) {
	args := m.Called(ctx, storeID)
	if o := args.Get(0); o != nil {
		return o.([]*domain.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderRepo) UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

type MockStoreRepo struct {
	mock.Mock
}

func (m *MockStoreRepo) Create(ctx context.Context, store *domain.Store) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

func (m *MockStoreRepo) GetByID(ctx context.Context, id int64) (*domain.Store, error) {
	args := m.Called(ctx, id)
	if s := args.Get(0); s != nil {
		return s.(*domain.Store), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStoreRepo) GetBySellerID(ctx context.Context, sellerID int64) (*domain.Store, error) {
	args := m.Called(ctx, sellerID)
	if s := args.Get(0); s != nil {
		return s.(*domain.Store), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStoreRepo) Update(ctx context.Context, store *domain.Store) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

func (m *MockStoreRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestOrderService_Create_EmptyItems(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockStoreRepo := new(MockStoreRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockStoreRepo)

	ctx := context.Background()
	order := &domain.Order{
		CustomerID: 1,
		StoreID:    10,
		Items:      []domain.OrderItem{},
	}

	err := orderService.Create(ctx, order)
	assert.ErrorIs(t, err, appErr.ErrEmptyOrder)
}

func TestOrderService_Create_InsufficientStock(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockStoreRepo := new(MockStoreRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockStoreRepo)

	ctx := context.Background()
	order := &domain.Order{
		CustomerID: 1,
		StoreID:    10,
		Items: []domain.OrderItem{
			{ProductID: 100, Quantity: 50, Price: 0.01}, // Fake client price
		},
	}

	mockOrderRepo.On("Create", ctx, order).Return(appErr.ErrInsufficientStock)

	err := orderService.Create(ctx, order)
	assert.ErrorIs(t, err, appErr.ErrInsufficientStock)
	mockOrderRepo.AssertExpectations(t)
}

func TestOrderService_Create_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockStoreRepo := new(MockStoreRepo)
	orderService := service.NewOrderService(mockOrderRepo, mockStoreRepo)

	ctx := context.Background()
	order := &domain.Order{
		CustomerID: 1,
		StoreID:    10,
		Items: []domain.OrderItem{
			{ProductID: 100, Quantity: 2, Price: 999.99}, // Client price will be ignored in DB repo
		},
	}

	mockOrderRepo.On("Create", ctx, order).Return(nil)

	err := orderService.Create(ctx, order)
	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
}
