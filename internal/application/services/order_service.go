package services

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/domain/repositories"
	"github.com/AndrivA89/orders/internal/domain/services"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo   repositories.OrderRepository
	userRepo    repositories.UserRepository
	productRepo repositories.ProductRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	productRepo repositories.ProductRepository,
) services.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, request *services.OrderRequest) (*entities.Order, error) {
	if _, err := s.userRepo.GetByID(ctx, request.UserID); err != nil {
		return nil, err
	}

	order := entities.NewOrder(request.UserID)

	for _, itemReq := range request.Items {
		product, err := s.productRepo.GetByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, err
		}

		if err := order.AddItem(product, itemReq.Quantity); err != nil {
			return nil, err
		}
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
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
