package entities

import (
	"time"

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
