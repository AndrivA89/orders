package entities

import (
	"time"

	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Quantity    int       `json:"quantity"`
	Price       int64     `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) ValidateForCreation() error {
	if p.Description == "" {
		return domainErrors.ErrProductDescriptionRequired
	}

	if p.Price <= 0 {
		return domainErrors.ErrProductPriceInvalid
	}

	if p.Quantity < 0 {
		return domainErrors.ErrProductQuantityNegative
	}

	return nil
}

func (p *Product) IsAvailable(requestedQuantity int) bool {
	return p.Quantity >= requestedQuantity
}

func (p *Product) ReserveQuantity(quantity int) error {
	if quantity <= 0 {
		return domainErrors.ErrQuantityInvalid
	}

	if !p.IsAvailable(quantity) {
		return domainErrors.ErrInsufficientQuantity
	}

	p.Quantity -= quantity

	return nil
}
