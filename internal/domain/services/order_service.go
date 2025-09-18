package services

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderRequest struct {
	UserID uuid.UUID
	Items  []OrderItemRequest
}

type OrderItemRequest struct {
	ProductID uuid.UUID
	Quantity  int
}

// TODO: CreateOrderWithTransaction - для production нужны транзакции

type OrderService interface {
	CreateOrder(ctx context.Context, request *OrderRequest) (*entities.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error)
	ConfirmOrder(ctx context.Context, orderID uuid.UUID) error
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
}
