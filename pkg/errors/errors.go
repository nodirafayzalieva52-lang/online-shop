package pkg

import "errors"

var (
	ErrEmptyOrder   = errors.New("order must contain at least one item")
	ErrInvalidOrder = errors.New("invalid order data")
)
