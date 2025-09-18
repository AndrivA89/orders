package services

import (
	"context"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req *services.CreateUserRequest) (*entities.User, error) {
	user := &entities.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		IsMarried: req.IsMarried,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.ValidateForCreation(req.Password); err != nil {
		return nil, err
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
