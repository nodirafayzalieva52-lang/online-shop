package dto

import (
	"time"

	"shop/internal/domain"
)

// APIError represents structured API error response payload.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents standard error response container.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// RegisterRequest payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest payload for user authentication.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse payload containing JWT token.
type LoginResponse struct {
	Token string `json:"token"`
}

// UserResponse payload for user profile endpoint.
type UserResponse struct {
	ID        int64       `json:"id"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// UpdateUserRequest payload for PATCH /me.
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

// CreateStoreRequest payload for POST /stores.
type CreateStoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateStoreRequest payload for PATCH /stores/{id}.
type UpdateStoreRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateCategoryRequest payload for POST /categories.
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// UpdateCategoryRequest payload for PATCH /categories/{id}.
type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

// CreateProductRequest payload for POST /products.
type CreateProductRequest struct {
	StoreID     int64   `json:"store_id"`
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProductRequest payload for PATCH /products/{id}.
type UpdateProductRequest struct {
	CategoryID  *int64   `json:"category_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

// OrderItemRequest payload within order creation request.
type OrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	StoreID   int64 `json:"store_id,omitempty"`
	Quantity  int   `json:"quantity"`
}

// CreateOrderRequest payload for POST /orders.
type CreateOrderRequest struct {
	StoreID int64              `json:"store_id"`
	Items   []OrderItemRequest `json:"items"`
}

// UpdateOrderStatusRequest payload for PATCH /orders/{id}/status.
type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status"`
}
