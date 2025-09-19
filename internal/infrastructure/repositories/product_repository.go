package repositories

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/infrastructure/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	model := &models.ProductModel{}
	if err := model.FromEntity(product); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	product.ID = model.ID
	product.CreatedAt = model.CreatedAt
	product.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var model models.ProductModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *productRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var model models.ProductModel
	// SELECT ... FOR UPDATE для предотвращения race conditions
	if err := r.db.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	var productModels []models.ProductModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&productModels).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Product, len(productModels))
	for i, model := range productModels {
		result[i] = model.ToEntity()
	}

	return result, nil
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	model := &models.ProductModel{}
	if err := model.FromEntity(product); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(model).Error
}
