package services

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*entities.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]*entities.Product, error)
}
