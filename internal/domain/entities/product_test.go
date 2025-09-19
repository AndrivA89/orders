package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProduct_ValidateForCreation(t *testing.T) {
	tests := []struct {
		name        string
		product     Product
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid product",
			product: Product{
				Description: "Valid Product",
				Price:       1000,
				Quantity:    5,
			},
			expectError: false,
		},
		{
			name: "empty description",
			product: Product{
				Description: "",
				Price:       1000,
				Quantity:    5,
			},
			expectError: true,
			errorMsg:    "description is required",
		},
		{
			name: "zero price",
			product: Product{
				Description: "Product",
				Price:       0,
				Quantity:    5,
			},
			expectError: true,
			errorMsg:    "price must be greater than 0",
		},
		{
			name: "negative price",
			product: Product{
				Description: "Product",
				Price:       -100,
				Quantity:    5,
			},
			expectError: true,
			errorMsg:    "price must be greater than 0",
		},
		{
			name: "negative quantity",
			product: Product{
				Description: "Product",
				Price:       1000,
				Quantity:    -1,
			},
			expectError: true,
			errorMsg:    "quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.ValidateForCreation()

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	product := Product{
		Quantity: 10,
	}

	assert.True(t, product.IsAvailable(5))
	assert.True(t, product.IsAvailable(10))
	assert.False(t, product.IsAvailable(15))
	assert.False(t, product.IsAvailable(11))
}

func TestProduct_ReserveQuantity(t *testing.T) {
	product := Product{
		ID:       uuid.New(),
		Quantity: 10,
	}

	// Successful reservation
	err := product.ReserveQuantity(5)
	assert.NoError(t, err)
	assert.Equal(t, 5, product.Quantity)

	// Try to reserve more than available
	err = product.ReserveQuantity(10)
	assert.Error(t, err)
	assert.Equal(t, "insufficient quantity available", err.Error())
	assert.Equal(t, 5, product.Quantity) // Quantity should not change

	// Try to reserve zero or negative
	err = product.ReserveQuantity(0)
	assert.Error(t, err)
	assert.Equal(t, "quantity must be greater than 0", err.Error())

	err = product.ReserveQuantity(-1)
	assert.Error(t, err)
	assert.Equal(t, "quantity must be greater than 0", err.Error())
}
