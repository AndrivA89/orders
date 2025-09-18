package repositories

//go:generate mockgen -source=order_repository.go -destination=mocks/order_repository_mock.go -package=mocks

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

// OrderRepository определяет контракт для работы с заказами
type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
}
