package repository

import (
	"context"

	"shop/internal/domain"
)

// UserRepository defines database operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateRole(ctx context.Context, userID int64, role domain.Role) error
}

// StoreRepository defines database operations for stores.
type StoreRepository interface {
	Create(ctx context.Context, store *domain.Store) error
	GetByID(ctx context.Context, id int64) (*domain.Store, error)
	GetBySellerID(ctx context.Context, sellerID int64) (*domain.Store, error)
	Update(ctx context.Context, store *domain.Store) error
	Delete(ctx context.Context, id int64) error
}

// CategoryRepository defines database operations for categories.
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetAll(ctx context.Context) ([]*domain.Category, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
}

// ProductRepository defines database operations for products.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	GetByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	UpdateStock(ctx context.Context, productID int64, delta int) error
}

// OrderRepository defines database operations for orders.
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
	GetByCustomerID(ctx context.Context, customerID int64) ([]*domain.Order, error)
	GetByStoreID(ctx context.Context, storeID int64) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, newStatus domain.OrderStatus) error
}
