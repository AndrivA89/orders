package repositories

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repositories.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
	model := &models.OrderModel{}
	if err := model.FromEntity(order); err != nil {
		return err
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(model).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	order.ID = model.ID
	order.CreatedAt = model.CreatedAt
	order.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var model models.OrderModel
	if err := r.db.WithContext(ctx).Preload("Items").First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity()
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	var models []models.OrderModel
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Find(&models).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Order, len(models))
	for i, model := range models {
		order, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		result[i] = order
	}

	return result, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) error {
	model := &models.OrderModel{}
	if err := model.FromEntity(order); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.OrderModel{}, "id = ?", id).Error
}
