package service_test

import (
	"context"
	"testing"

	"shop/internal/delivery/http/dto"
	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if p := args.Get(0); p != nil {
		return p.(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) GetByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, storeID, limit, offset)
	if p := args.Get(0); p != nil {
		return p.([]*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	args := m.Called(ctx, limit, offset)
	if p := args.Get(0); p != nil {
		return p.([]*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepo) UpdateStock(ctx context.Context, productID int64, delta int) error {
	args := m.Called(ctx, productID, delta)
	return args.Error(0)
}

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepo) GetAll(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	if c := args.Get(0); c != nil {
		return c.([]*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if c := args.Get(0); c != nil {
		return c.(*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestProductService_CreateProduct_InvalidInputs(t *testing.T) {
	mockProdRepo := new(MockProductRepo)
	mockStoreRepo := new(MockStoreRepo)
	svc := service.NewProductService(mockProdRepo, mockStoreRepo)
	ctx := context.Background()

	// Empty name
	prod, err := svc.CreateProduct(ctx, 1, domain.RoleSeller, dto.CreateProductRequest{
		StoreID: 1, Name: "", Price: 100, Stock: 5,
	})
	assert.Nil(t, prod)
	assert.ErrorContains(t, err, "cannot be empty")

	// Negative price
	prod, err = svc.CreateProduct(ctx, 1, domain.RoleSeller, dto.CreateProductRequest{
		StoreID: 1, Name: "Item", Price: -10, Stock: 5,
	})
	assert.Nil(t, prod)
	assert.ErrorContains(t, err, "greater than zero")

	// Negative stock
	prod, err = svc.CreateProduct(ctx, 1, domain.RoleSeller, dto.CreateProductRequest{
		StoreID: 1, Name: "Item", Price: 100, Stock: -1,
	})
	assert.Nil(t, prod)
	assert.ErrorContains(t, err, "cannot be negative")
}

func TestProductService_CreateProduct_AccessDenied(t *testing.T) {
	mockProdRepo := new(MockProductRepo)
	mockStoreRepo := new(MockStoreRepo)
	svc := service.NewProductService(mockProdRepo, mockStoreRepo)
	ctx := context.Background()

	store := &domain.Store{ID: 1, SellerID: 99, Name: "Other Store"}
	mockStoreRepo.On("GetByID", ctx, int64(1)).Return(store, nil)

	// User 10 trying to add product to Store owned by User 99
	prod, err := svc.CreateProduct(ctx, 10, domain.RoleSeller, dto.CreateProductRequest{
		StoreID: 1, Name: "Unauthorized Item", Price: 100, Stock: 5,
	})

	assert.Nil(t, prod)
	assert.ErrorIs(t, err, appErr.ErrAccessDenied)
	mockStoreRepo.AssertExpectations(t)
}

func TestStoreService_CreateStore_InvalidInputs(t *testing.T) {
	mockStoreRepo := new(MockStoreRepo)
	mockUserRepo := new(MockUserRepo)
	svc := service.NewStoreService(mockStoreRepo, mockUserRepo)
	ctx := context.Background()

	// Empty store name
	store, err := svc.CreateStore(ctx, 1, "", "Description")
	assert.Nil(t, store)
	assert.ErrorContains(t, err, "cannot be empty")
}

func TestCategoryService_CreateCategory_InvalidInputs(t *testing.T) {
	mockCatRepo := new(MockCategoryRepo)
	svc := service.NewCategoryService(mockCatRepo)
	ctx := context.Background()

	// Empty category name
	cat, err := svc.CreateCategory(ctx, "")
	assert.Nil(t, cat)
	assert.ErrorContains(t, err, "cannot be empty")
}
