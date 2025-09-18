package models

import (
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FirstName string         `gorm:"column:first_name;not null;size:100" json:"first_name"`
	LastName  string         `gorm:"column:last_name;not null;size:100" json:"last_name"`
	Age       int            `gorm:"column:age;not null" json:"age"`
	IsMarried bool           `gorm:"column:is_married;default:false" json:"is_married"`
	Password  string         `gorm:"column:password;not null;size:255" json:"-"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Orders []OrderModel `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

func (u *UserModel) ToEntity() *entities.User {
	return &entities.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Age:       u.Age,
		IsMarried: u.IsMarried,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *UserModel) FromEntity(entity *entities.User) {
	u.ID = entity.ID
	u.FirstName = entity.FirstName
	u.LastName = entity.LastName
	u.Age = entity.Age
	u.IsMarried = entity.IsMarried
	u.Password = entity.Password
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
}
