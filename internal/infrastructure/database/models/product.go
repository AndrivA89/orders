package models

import (
	"encoding/json"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Description string         `gorm:"column:description;not null;size:500" json:"description"`
	Tags        datatypes.JSON `gorm:"column:tags;type:json" json:"tags"`
	Quantity    int            `gorm:"column:quantity;not null;default:0" json:"quantity"`
	Price       int64          `gorm:"column:price;not null" json:"price"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (p *ProductModel) ToEntity() *entities.Product {
	var tags []string
	if p.Tags != nil {
		// Игнорируем ошибку - если JSON невалидный, используем пустой слайс
		_ = json.Unmarshal(p.Tags, &tags)
	}

	return &entities.Product{
		ID:          p.ID,
		Description: p.Description,
		Tags:        tags,
		Quantity:    p.Quantity,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (p *ProductModel) FromEntity(entity *entities.Product) error {
	p.ID = entity.ID
	p.Description = entity.Description
	p.Quantity = entity.Quantity
	p.Price = entity.Price
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt

	if entity.Tags != nil {
		tags, err := json.Marshal(entity.Tags)
		if err != nil {
			return err
		}
		p.Tags = tags
	}

	return nil
}
