package dto

import (
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID uuid.UUID          `json:"user_id" binding:"required"`
	Items  []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	ID        uuid.UUID           `json:"id"`
	UserID    uuid.UUID           `json:"user_id"`
	Status    string              `json:"status"`
	Total     int64               `json:"total"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID              uuid.UUID               `json:"id"`
	ProductID       uuid.UUID               `json:"product_id"`
	ProductSnapshot ProductSnapshotResponse `json:"product_snapshot"`
	Quantity        int                     `json:"quantity"`
	PricePerItem    int64                   `json:"price_per_item"`
	Total           int64                   `json:"total"`
	CreatedAt       time.Time               `json:"created_at"`
}

type ProductSnapshotResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Price       int64     `json:"price"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

func ToOrderResponse(order *entities.Order) *OrderResponse {
	items := make([]OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			ProductSnapshot: ProductSnapshotResponse{
				ID:          item.ProductSnapshot.ID,
				Description: item.ProductSnapshot.Description,
				Tags:        item.ProductSnapshot.Tags,
				Price:       item.ProductSnapshot.Price,
			},
			Quantity:     item.Quantity,
			PricePerItem: item.PricePerItem,
			Total:        item.Total,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    string(order.Status),
		Total:     order.Total,
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func ToOrderListResponse(orders []*entities.Order, limit, offset int) *OrderListResponse {
	orderResponses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, *ToOrderResponse(order))
	}

	return &OrderListResponse{
		Orders: orderResponses,
		Total:  len(orderResponses),
		Limit:  limit,
		Offset: offset,
	}
}
