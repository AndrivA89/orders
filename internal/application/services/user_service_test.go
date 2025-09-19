package services

import (
	"context"
	"testing"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories/mocks"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_RegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockUserRepo)

	request := &services.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
		IsMarried: false,
		Password:  "password123",
	}

	mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, user *entities.User) error {
			user.ID = uuid.New()
			return nil
		})

	user, err := service.RegisterUser(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, 25, user.Age)
	assert.False(t, user.IsMarried)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "password123", user.Password) // Password should be hashed
}

func TestUserService_RegisterUser_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockUserRepo)

	request := &services.CreateUserRequest{
		FirstName: "", // Empty first name
		LastName:  "Doe",
		Age:       25,
		Password:  "password123",
	}

	user, err := service.RegisterUser(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "first name is required", err.Error())
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockUserRepo)

	userID := uuid.New()
	expectedUser := &entities.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
	}

	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(expectedUser, nil)

	user, err := service.GetUserByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "John", user.FirstName)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(mockUserRepo)

	userID := uuid.New()

	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, assert.AnError)

	user, err := service.GetUserByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, user)
}
