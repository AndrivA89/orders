package entities

import (
	"strings"
	"time"

	"github.com/AndrivA89/orders/internal/domain/constants"
	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"age"`
	IsMarried bool      `json:"is_married"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetFullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

func (u *User) ValidateForCreation(plainPassword string) error {
	if u.FirstName == "" {
		return domainErrors.ErrFirstNameRequired
	}

	if u.LastName == "" {
		return domainErrors.ErrLastNameRequired
	}

	if u.Age < constants.MinUserAge {
		return domainErrors.ErrUserTooYoung
	}

	if len(plainPassword) < constants.MinPasswordLength {
		return domainErrors.ErrPasswordTooShort
	}

	return nil
}

func (u *User) SetPassword(plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), constants.BcryptCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))

	return err == nil
}
