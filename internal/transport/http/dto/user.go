package dto

import (
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
	Age       int    `json:"age" binding:"required,min=18"`
	IsMarried *bool  `json:"is_married"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Age       int       `json:"age"`
	IsMarried bool      `json:"is_married"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),
		Age:       user.Age,
		IsMarried: user.IsMarried,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
