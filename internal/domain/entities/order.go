package entities

import (
	"time"

	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	Status    OrderStatus `json:"status"`
	Total     int64       `json:"total"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID              uuid.UUID       `json:"id"`
	OrderID         uuid.UUID       `json:"order_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	ProductSnapshot ProductSnapshot `json:"product_snapshot"`
	Quantity        int             `json:"quantity"`
	PricePerItem    int64           `json:"price_per_item"`
	Total           int64           `json:"total"`
	CreatedAt       time.Time       `json:"created_at"`
}

type ProductSnapshot struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Price       int64     `json:"price"`
}

func NewOrder(userID uuid.UUID) *Order {
	return &Order{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    OrderStatusPending,
		Items:     make([]OrderItem, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (o *Order) AddItem(product *Product, quantity int) error {
	if quantity <= 0 {
		return domainErrors.ErrQuantityInvalid
	}

	if !product.IsAvailable(quantity) {
		return domainErrors.ErrInsufficientStock
	}

	snapshot := ProductSnapshot{
		ID:          product.ID,
		Description: product.Description,
		Tags:        product.Tags,
		Price:       product.Price,
	}

	item := OrderItem{
		ID:              uuid.New(),
		OrderID:         o.ID,
		ProductID:       product.ID,
		ProductSnapshot: snapshot,
		Quantity:        quantity,
		PricePerItem:    product.Price,
		Total:           product.Price * int64(quantity),
		CreatedAt:       time.Now(),
	}

	o.Items = append(o.Items, item)
	o.calculateTotal()
	o.UpdatedAt = time.Now()

	return nil
}

func (o *Order) calculateTotal() {
	var total int64

	for _, item := range o.Items {
		total += item.Total
	}

	o.Total = total
}

func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return domainErrors.ErrOnlyPendingCanConfirm
	}
	if len(o.Items) == 0 {
		return domainErrors.ErrCannotConfirmEmptyOrder
	}

	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()

	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusCompleted {
		return domainErrors.ErrCompletedOrdersReadonly
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()

	return nil
}
