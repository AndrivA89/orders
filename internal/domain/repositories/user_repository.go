package repositories

//go:generate mockgen -source=user_repository.go -destination=mocks/user_repository_mock.go -package=mocks

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository определяет контракт для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}
