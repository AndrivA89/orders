package models

import (
	"encoding/json"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Status    string         `gorm:"column:status;not null;size:20;default:'pending'" json:"status"`
	Total     int64          `gorm:"column:total;not null;default:0" json:"total"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User  UserModel        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []OrderItemModel `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItemModel struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductSnapshot datatypes.JSON `gorm:"column:product_snapshot;type:json;not null" json:"product_snapshot"`
	Quantity        int            `gorm:"column:quantity;not null" json:"quantity"`
	PricePerItem    int64          `gorm:"column:price_per_item;not null" json:"price_per_item"`
	Total           int64          `gorm:"column:total;not null" json:"total"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Order   OrderModel   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product ProductModel `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (o *OrderModel) ToEntity() (*entities.Order, error) {
	order := &entities.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		Status:    entities.OrderStatus(o.Status),
		Total:     o.Total,
		Items:     make([]entities.OrderItem, 0, len(o.Items)),
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}

	for _, item := range o.Items {
		orderItem, err := item.ToEntity()
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, *orderItem)
	}

	return order, nil
}

func (o *OrderModel) FromEntity(entity *entities.Order) error {
	o.ID = entity.ID
	o.UserID = entity.UserID
	o.Status = string(entity.Status)
	o.Total = entity.Total
	o.CreatedAt = entity.CreatedAt
	o.UpdatedAt = entity.UpdatedAt

	o.Items = make([]OrderItemModel, 0, len(entity.Items))
	for _, item := range entity.Items {
		itemModel := &OrderItemModel{}
		if err := itemModel.FromEntity(&item); err != nil {
			return err
		}
		o.Items = append(o.Items, *itemModel)
	}

	return nil
}

func (oi *OrderItemModel) ToEntity() (*entities.OrderItem, error) {
	var snapshot entities.ProductSnapshot
	if err := json.Unmarshal(oi.ProductSnapshot, &snapshot); err != nil {
		return nil, err
	}

	return &entities.OrderItem{
		ID:              oi.ID,
		OrderID:         oi.OrderID,
		ProductID:       oi.ProductID,
		ProductSnapshot: snapshot,
		Quantity:        oi.Quantity,
		PricePerItem:    oi.PricePerItem,
		Total:           oi.Total,
		CreatedAt:       oi.CreatedAt,
	}, nil
}

func (oi *OrderItemModel) FromEntity(entity *entities.OrderItem) error {
	oi.ID = entity.ID
	oi.OrderID = entity.OrderID
	oi.ProductID = entity.ProductID
	oi.Quantity = entity.Quantity
	oi.PricePerItem = entity.PricePerItem
	oi.Total = entity.Total
	oi.CreatedAt = entity.CreatedAt

	snapshot, err := json.Marshal(entity.ProductSnapshot)
	if err != nil {
		return err
	}
	oi.ProductSnapshot = snapshot

	return nil
}
