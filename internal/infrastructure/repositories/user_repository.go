package repositories

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	model := &models.UserModel{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}
