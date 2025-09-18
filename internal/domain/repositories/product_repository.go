package repositories

//go:generate mockgen -source=product_repository.go -destination=mocks/product_repository_mock.go -package=mocks

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

// ProductRepository определяет контракт для работы с товарами
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
}
