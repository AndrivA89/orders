package services

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *CreateUserRequest) (*entities.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}
