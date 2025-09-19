package services

import (
	"context"
	"testing"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories/mocks"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductService_CreateProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := NewProductService(mockProductRepo)

	request := &services.CreateProductRequest{
		Description: "Test Product",
		Price:       1000,
		Quantity:    10,
	}

	mockProductRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, product *entities.Product) error {
			product.ID = uuid.New()
			product.CreatedAt = time.Now()
			product.UpdatedAt = time.Now()
			return nil
		})

	product, err := service.CreateProduct(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Description)
	assert.Equal(t, int64(1000), product.Price)
	assert.Equal(t, 10, product.Quantity)
}

func TestProductService_CreateProduct_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := NewProductService(mockProductRepo)

	request := &services.CreateProductRequest{
		Description: "", // Empty description
		Price:       1000,
		Quantity:    10,
	}

	product, err := service.CreateProduct(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "description is required", err.Error())
}

func TestProductService_GetProductByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := NewProductService(mockProductRepo)

	productID := uuid.New()
	expectedProduct := &entities.Product{
		ID:          productID,
		Description: "Test Product",
		Price:       1000,
		Quantity:    10,
	}

	mockProductRepo.EXPECT().GetByID(gomock.Any(), productID).Return(expectedProduct, nil)

	product, err := service.GetProductByID(context.Background(), productID)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, productID, product.ID)
	assert.Equal(t, "Test Product", product.Description)
}

func TestProductService_GetProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	service := NewProductService(mockProductRepo)

	expectedProducts := []*entities.Product{
		{
			ID:          uuid.New(),
			Description: "Product 1",
			Price:       1000,
			Quantity:    5,
		},
		{
			ID:          uuid.New(),
			Description: "Product 2",
			Price:       2000,
			Quantity:    3,
		},
	}

	mockProductRepo.EXPECT().GetAll(gomock.Any(), 10, 0).Return(expectedProducts, nil)

	products, err := service.GetProducts(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, products)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Description)
	assert.Equal(t, "Product 2", products[1].Description)
}
