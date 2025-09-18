package handlers

import (
	"net/http"
	"strconv"

	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"
	"github.com/AndrivA89/orders/internal/domain/services"
	"github.com/AndrivA89/orders/internal/transport/http/dto"
	"github.com/AndrivA89/orders/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), req.ToServiceRequest())
	if err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToOrderResponse(order))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.HandleValidationError(c, domainErrors.ErrInvalidOrderID)
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		middleware.HandleNotFoundError(c, domainErrors.ErrOrderNotFound)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(order))
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		middleware.HandleValidationError(c, domainErrors.ErrInvalidUserID)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	orders, err := h.orderService.GetOrdersByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		middleware.HandleInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderListResponse(orders, limit, offset))
}

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.HandleValidationError(c, domainErrors.ErrInvalidOrderID)
		return
	}

	if err := h.orderService.ConfirmOrder(c.Request.Context(), orderID); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		middleware.HandleInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(order))
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.HandleValidationError(c, domainErrors.ErrInvalidOrderID)
		return
	}

	if err := h.orderService.CancelOrder(c.Request.Context(), orderID); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		middleware.HandleInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToOrderResponse(order))
}
