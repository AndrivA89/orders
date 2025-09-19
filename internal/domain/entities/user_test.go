package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_ValidateForCreation(t *testing.T) {
	tests := []struct {
		name        string
		user        User
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid user",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
			},
			password:    "password123",
			expectError: false,
		},
		{
			name: "empty first name",
			user: User{
				FirstName: "",
				LastName:  "Doe",
				Age:       25,
			},
			password:    "password123",
			expectError: true,
			errorMsg:    "first name is required",
		},
		{
			name: "empty last name",
			user: User{
				FirstName: "John",
				LastName:  "",
				Age:       25,
			},
			password:    "password123",
			expectError: true,
			errorMsg:    "last name is required",
		},
		{
			name: "age less than 18",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       17,
			},
			password:    "password123",
			expectError: true,
			errorMsg:    "user must be at least 18 years old",
		},
		{
			name: "password too short",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Age:       25,
			},
			password:    "short",
			expectError: true,
			errorMsg:    "password must be at least 8 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateForCreation(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_GetFullName(t *testing.T) {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	fullName := user.GetFullName()
	assert.Equal(t, "John Doe", fullName)
}

func TestUser_SetPassword(t *testing.T) {
	user := User{}
	password := "testpassword123"

	err := user.SetPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password) // Пароль должен быть захеширован
}

func TestUser_CheckPassword(t *testing.T) {
	user := User{
		ID: uuid.New(),
	}
	password := "testpassword123"

	err := user.SetPassword(password)
	assert.NoError(t, err)

	assert.True(t, user.CheckPassword(password))
	assert.False(t, user.CheckPassword("wrongpassword"))
}
