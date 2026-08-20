package errors

import "errors"

var (
	ErrEmptyOrder              = errors.New("order must contain at least one item")
	ErrInvalidOrder            = errors.New("invalid order data")
	ErrInvalidToken            = errors.New("invalid or expired token")
	ErrOrderNotFound           = errors.New("order not found")
	ErrAccessDenied            = errors.New("access denied")
	ErrNotFound                = errors.New("resource not found")
	ErrAlreadyExists           = errors.New("resource already exists")
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrUserNotFound            = errors.New("user not found")
	ErrStoreNotFound           = errors.New("store not found")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrProductNotFound         = errors.New("product not found")
	ErrInsufficientStock       = errors.New("insufficient stock for product")
	ErrInvalidStore            = errors.New("product does not belong to specified store")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
)
