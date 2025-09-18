package dto

import (
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Description string   `json:"description" binding:"required,min=1,max=500"`
	Tags        []string `json:"tags"`
	Quantity    int      `json:"quantity" binding:"required,min=0"`
	Price       int64    `json:"price" binding:"required,min=1"`
}

func (req *CreateProductRequest) ToServiceRequest() *services.CreateProductRequest {
	return &services.CreateProductRequest{
		Description: req.Description,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
		Price:       req.Price,
	}
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Quantity    int       `json:"quantity"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToProductResponse(product *entities.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		Description: product.Description,
		Tags:        product.Tags,
		Quantity:    product.Quantity,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
