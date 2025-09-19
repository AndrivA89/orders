package services

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"
	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo   repositories.OrderRepository
	userRepo    repositories.UserRepository
	productRepo repositories.ProductRepository
	txManager   repositories.TransactionManager
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	productRepo repositories.ProductRepository,
	txManager repositories.TransactionManager,
) services.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, request *services.OrderRequest) (*entities.Order, error) {
	if len(request.Items) == 0 {
		return nil, domainErrors.ErrOrderMustHaveItems
	}

	var resultOrder *entities.Order

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, repos repositories.TransactionalRepositories) error {
		// Verify user exists
		_, err := repos.UserRepository.GetByID(ctx, request.UserID)
		if err != nil {
			return domainErrors.ErrUserNotFound
		}

		order := entities.NewOrder(request.UserID)

		// Process each item with quantity reservation
		for _, itemReq := range request.Items {
			// Lock product row to prevent race conditions
			product, err := repos.ProductRepository.GetByIDForUpdate(ctx, itemReq.ProductID)
			if err != nil {
				return err
			}

			// Add item to order (includes availability check)
			if err := order.AddItem(product, itemReq.Quantity); err != nil {
				return err
			}

			// Reserve quantity in product
			if err := product.ReserveQuantity(itemReq.Quantity); err != nil {
				return err
			}

			// Update product with reserved quantity
			if err := repos.ProductRepository.Update(ctx, product); err != nil {
				return err
			}
		}

		// Create order with all items
		if err := repos.OrderRepository.Create(ctx, order); err != nil {
			return err
		}

		resultOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultOrder, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *orderService) GetOrdersByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *orderService) ConfirmOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.Confirm(); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, order)
}

func (s *orderService) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, order)
}
