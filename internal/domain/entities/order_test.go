package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, OrderStatusPending, order.Status)
	assert.Equal(t, int64(0), order.Total)
	assert.Empty(t, order.Items)
	assert.NotEqual(t, uuid.Nil, order.ID)
}

func TestOrder_AddItem(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	product := &Product{
		ID:          uuid.New(),
		Description: "Test Product",
		Quantity:    10,
		Price:       1000, // 10.00 в копейках
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Успешное добавление товара
	err := order.AddItem(product, 2)
	assert.NoError(t, err)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, int64(2000), order.Total) // 2 * 1000

	item := order.Items[0]
	assert.Equal(t, product.ID, item.ProductID)
	assert.Equal(t, 2, item.Quantity)
	assert.Equal(t, product.Price, item.PricePerItem)
	assert.Equal(t, int64(2000), item.Total)
	assert.Equal(t, product.Description, item.ProductSnapshot.Description)
}

func TestOrder_AddItem_InvalidQuantity(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	product := &Product{
		ID:       uuid.New(),
		Quantity: 10,
		Price:    1000,
	}

	err := order.AddItem(product, 0)
	assert.Error(t, err)
	assert.Equal(t, "quantity must be greater than 0", err.Error())

	err = order.AddItem(product, -1)
	assert.Error(t, err)
}

func TestOrder_AddItem_InsufficientQuantity(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	product := &Product{
		ID:       uuid.New(),
		Quantity: 5,
		Price:    1000,
	}

	err := order.AddItem(product, 10)
	assert.Error(t, err)
	assert.Equal(t, "insufficient product quantity", err.Error())
}

func TestOrder_Confirm(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	product := &Product{
		ID:       uuid.New(),
		Quantity: 10,
		Price:    1000,
	}
	err := order.AddItem(product, 2)
	assert.NoError(t, err)

	err = order.Confirm()
	assert.NoError(t, err)
	assert.Equal(t, OrderStatusConfirmed, order.Status)
}

func TestOrder_Confirm_EmptyOrder(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	err := order.Confirm()
	assert.Error(t, err)
	assert.Equal(t, "cannot confirm empty order", err.Error())
}

func TestOrder_Cancel(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)

	err := order.Cancel()
	assert.NoError(t, err)
	assert.Equal(t, OrderStatusCancelled, order.Status)
}

func TestOrder_Cancel_CompletedOrder(t *testing.T) {
	userID := uuid.New()
	order := NewOrder(userID)
	order.Status = OrderStatusCompleted

	err := order.Cancel()
	assert.Error(t, err)
	assert.Equal(t, "completed orders cannot be cancelled", err.Error())
}
