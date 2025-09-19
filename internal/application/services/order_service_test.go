package services

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/domain/repositories/mocks"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestOrderService_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	userID := uuid.New()
	productID := uuid.New()

	user := &entities.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
	}

	product := &entities.Product{
		ID:          productID,
		Description: "Test Product",
		Quantity:    10,
		Price:       1000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	request := &services.OrderRequest{
		UserID: userID,
		Items: []services.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockTxManager.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context, repositories.TransactionalRepositories) error) error {
			repos := repositories.TransactionalRepositories{
				OrderRepository:   mockOrderRepo,
				ProductRepository: mockProductRepo,
				UserRepository:    mockUserRepo,
			}
			return fn(ctx, repos)
		},
	)
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
	mockProductRepo.EXPECT().GetByIDForUpdate(gomock.Any(), productID).Return(product, nil)
	mockProductRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockOrderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	order, err := service.CreateOrder(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, int64(2000), order.Total)
}

func TestOrderService_CreateOrder_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)
	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	userID := uuid.New()
	productID := uuid.New()

	request := &services.OrderRequest{
		UserID: userID,
		Items: []services.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockTxManager.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context, repositories.TransactionalRepositories) error) error {
			repos := repositories.TransactionalRepositories{
				OrderRepository:   mockOrderRepo,
				ProductRepository: mockProductRepo,
				UserRepository:    mockUserRepo,
			}
			return fn(ctx, repos)
		},
	)
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, gorm.ErrRecordNotFound)

	order, err := service.CreateOrder(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "user not found", err.Error())
}

func TestOrderService_CreateOrder_InsufficientQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)
	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	userID := uuid.New()
	productID := uuid.New()

	user := &entities.User{ID: userID}
	product := &entities.Product{
		ID:       productID,
		Quantity: 1, // Недостаточно товара
		Price:    1000,
	}

	request := &services.OrderRequest{
		UserID: userID,
		Items: []services.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  5, // Заказываем больше чем есть
			},
		},
	}

	mockTxManager.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context, repositories.TransactionalRepositories) error) error {
			repos := repositories.TransactionalRepositories{
				OrderRepository:   mockOrderRepo,
				ProductRepository: mockProductRepo,
				UserRepository:    mockUserRepo,
			}
			return fn(ctx, repos)
		},
	)
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
	mockProductRepo.EXPECT().GetByIDForUpdate(gomock.Any(), productID).Return(product, nil)

	order, err := service.CreateOrder(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "insufficient product quantity", err.Error())
}

func TestOrderService_CreateOrder_EmptyItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)
	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	userID := uuid.New()
	request := &services.OrderRequest{
		UserID: userID,
		Items:  []services.OrderItemRequest{}, // Пустой список товаров
	}

	order, err := service.CreateOrder(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, "order must contain at least one item", err.Error())
}

func TestOrderService_ConfirmOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)
	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	orderID := uuid.New()
	order := &entities.Order{
		ID:     orderID,
		Status: entities.OrderStatusPending,
		Items: []entities.OrderItem{
			{
				ProductID: uuid.New(),
				Quantity:  2,
			},
		},
	}

	mockOrderRepo.EXPECT().GetByID(gomock.Any(), orderID).Return(order, nil)
	mockOrderRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	err := service.ConfirmOrder(context.Background(), orderID)

	assert.NoError(t, err)
}

// TestOrderService_CreateOrder_RaceCondition проверяет корректность работы при конкурентном доступе
func TestOrderService_CreateOrder_RaceCondition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductRepo := mocks.NewMockProductRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTxManager := mocks.NewMockTransactionManager(ctrl)

	service := NewOrderService(mockOrderRepo, mockUserRepo, mockProductRepo, mockTxManager)

	userID1 := uuid.New()
	userID2 := uuid.New()
	productID := uuid.New()

	user1 := &entities.User{ID: userID1, FirstName: "User1", LastName: "Test"}
	user2 := &entities.User{ID: userID2, FirstName: "User2", LastName: "Test"}

	product := &entities.Product{
		ID:       productID,
		Quantity: 1,
		Price:    1000,
	}

	request1 := &services.OrderRequest{
		UserID: userID1,
		Items:  []services.OrderItemRequest{{ProductID: productID, Quantity: 1}},
	}

	request2 := &services.OrderRequest{
		UserID: userID2,
		Items:  []services.OrderItemRequest{{ProductID: productID, Quantity: 1}},
	}

	var (
		successful int
		failed     int
		wg         sync.WaitGroup
	)

	// Симулируем блокировку: первый запрос проходит, второй получает ошибку блокировки
	mockTxManager.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context, repositories.TransactionalRepositories) error) error {
			repos := repositories.TransactionalRepositories{
				OrderRepository:   mockOrderRepo,
				ProductRepository: mockProductRepo,
				UserRepository:    mockUserRepo,
			}
			return fn(ctx, repos)
		},
	).Times(2)

	// Настраиваем моки для первого успешного запроса
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID1).Return(user1, nil)
	mockProductRepo.EXPECT().GetByIDForUpdate(gomock.Any(), productID).Return(product, nil)
	mockProductRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, p *entities.Product) error {
			// Симулируем обновление количества
			p.Quantity = 0
			return nil
		})
	mockOrderRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// Настраиваем моки для второго неуспешного запроса
	mockUserRepo.EXPECT().GetByID(gomock.Any(), userID2).Return(user2, nil)
	mockProductRepo.EXPECT().GetByIDForUpdate(gomock.Any(), productID).DoAndReturn(
		func(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
			// Возвращаем товар с нулевым остатком (уже купил первый пользователь)
			return &entities.Product{
				ID:       productID,
				Quantity: 0,
				Price:    1000,
			}, nil
		})

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := service.CreateOrder(context.Background(), request1)
		if err == nil {
			successful++
		} else {
			failed++
		}
	}()

	go func() {
		defer wg.Done()
		_, err := service.CreateOrder(context.Background(), request2)
		if err == nil {
			successful++
		} else {
			failed++
		}
	}()

	wg.Wait()

	// Проверяем что только один заказ прошел
	assert.Equal(t, 1, successful, "Должен пройти только один заказ")
	assert.Equal(t, 1, failed, "Один заказ должен завершиться ошибкой")
}
