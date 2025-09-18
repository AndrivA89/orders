package handlers

import (
	"net/http"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	users []entities.User
}

func (h *UserHandler) GetUsers() []entities.User {
	return h.users
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: make([]entities.User, 0),
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isMarried := false
	if req.IsMarried != nil {
		isMarried = *req.IsMarried
	}

	user := entities.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		IsMarried: isMarried,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h.users = append(h.users, user)

	c.JSON(http.StatusCreated, dto.ToUserResponse(&user))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	for _, user := range h.users {
		if user.ID == userID {
			c.JSON(http.StatusOK, dto.ToUserResponse(&user))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
