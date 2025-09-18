package handlers

import (
	"net/http"

	"github.com/AndrivA89/orders/internal/domain/services"
	"github.com/AndrivA89/orders/internal/transport/http/dto"
	"github.com/AndrivA89/orders/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), req.ToServiceRequest())
	if err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		middleware.HandleNotFoundError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}
