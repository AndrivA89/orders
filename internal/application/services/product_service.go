package services

import (
	"context"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
)

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) services.ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *services.CreateProductRequest) (*entities.Product, error) {
	product := &entities.Product{
		ID:          uuid.New(),
		Description: req.Description,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := product.ValidateForCreation(); err != nil {
		return nil, err
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) GetProducts(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	return s.productRepo.GetAll(ctx, limit, offset)
}
