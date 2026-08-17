package repository

import (
	"context"

	"shop/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type StoreRepository interface {
	Create(ctx context.Context, store models.Store) error
	GetByID(ctx context.Context, id int) (*models.Store, error)
	GetBySellerID(ctx context.Context, sellerID int) (*models.Store, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category models.Category) error
	GetAll(ctx context.Context) ([]*models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
}

type ProductRespository interface {
	Create(ctx context.Context, product models.Product)  error
	GetByID(ctx context.Context, id int) (*models.Product, error)
	GettByStoreID(ctx context.Context, storeID int) ([]*models.Product, error)
	GetAll(ctx context.Context) ([]*models.Product, error)
	UpdateStock(ctx context.Context, productID int, delta int) error
}

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int) (*models.Order, error)
	GetByCustomerID(ctx context.Context, customersID int) ([]*models.Order, error)
	GetByStoreID(ctx context.Context, storeID int) ([]*models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}
