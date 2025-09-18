package errors

import "errors"

// User domain errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// Product domain errors
var (
	ErrProductDescriptionRequired = errors.New("description is required")
	ErrProductPriceInvalid        = errors.New("price must be greater than 0")
	ErrProductQuantityNegative    = errors.New("quantity cannot be negative")
	ErrInsufficientQuantity       = errors.New("insufficient quantity available")
	ErrProductNotFound            = errors.New("product not found")
)

// Order domain errors
var (
	ErrQuantityInvalid         = errors.New("quantity must be greater than 0")
	ErrInsufficientStock       = errors.New("insufficient product quantity")
	ErrOnlyPendingCanConfirm   = errors.New("only pending orders can be confirmed")
	ErrCannotConfirmEmptyOrder = errors.New("cannot confirm empty order")
	ErrCompletedOrdersReadonly = errors.New("completed orders cannot be cancelled")
	ErrOrderMustHaveItems      = errors.New("order must contain at least one item")
	ErrOrderNotFound           = errors.New("order not found")
)

// Validation errors
var (
	ErrInvalidOrderID = errors.New("invalid order ID format")
	ErrInvalidUserID  = errors.New("invalid user ID format")
)
