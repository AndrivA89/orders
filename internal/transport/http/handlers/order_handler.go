package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndrivA89/orders/internal/domain/entities"
	domainErrors "github.com/AndrivA89/orders/internal/domain/errors"
	"github.com/AndrivA89/orders/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orders   []entities.Order
	products []entities.Product
	users    []entities.User
}

func NewOrderHandler(productHandler *ProductHandler, userHandler *UserHandler) *OrderHandler {
	return &OrderHandler{
		orders:   make([]entities.Order, 0),
		products: productHandler.GetProductsData(),
		users:    userHandler.GetUsers(),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if !h.userExists(req.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrUserNotFound.Error()})
		return
	}

	order := entities.NewOrder(req.UserID)

	// Add items to order
	for _, itemReq := range req.Items {
		product := h.findProductByID(itemReq.ProductID)
		if product == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrProductNotFound.Error()})
			return
		}

		if err := order.AddItem(product, itemReq.Quantity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if len(order.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrOrderMustHaveItems.Error()})
		return
	}

	h.orders = append(h.orders, *order)

	c.JSON(http.StatusCreated, dto.ToOrderResponse(order))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrInvalidOrderID.Error()})
		return
	}

	for _, order := range h.orders {
		if order.ID == orderID {
			c.JSON(http.StatusOK, dto.ToOrderResponse(&order))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": domainErrors.ErrOrderNotFound.Error()})
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrInvalidUserID.Error()})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Find orders by user
	var userOrders []*entities.Order
	for i, order := range h.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, &h.orders[i])
		}
	}

	total := len(userOrders)
	orders := userOrders

	if offset >= total {
		orders = []*entities.Order{}
	} else {
		end := offset + limit
		if end > total {
			end = total
		}
		orders = userOrders[offset:end]
	}

	c.JSON(http.StatusOK, dto.ToOrderListResponse(orders, limit, offset))
}

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrInvalidOrderID.Error()})
		return
	}

	for i, order := range h.orders {
		if order.ID == orderID {
			if err := order.Confirm(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			h.orders[i] = order
			c.JSON(http.StatusOK, dto.ToOrderResponse(&order))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": domainErrors.ErrOrderNotFound.Error()})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrInvalidOrderID.Error()})
		return
	}

	for i, order := range h.orders {
		if order.ID == orderID {
			if err := order.Cancel(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			h.orders[i] = order
			c.JSON(http.StatusOK, dto.ToOrderResponse(&order))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": domainErrors.ErrOrderNotFound.Error()})
}

func (h *OrderHandler) userExists(userID uuid.UUID) bool {
	for _, user := range h.users {
		if user.ID == userID {
			return true
		}
	}
	return false
}

func (h *OrderHandler) findProductByID(productID uuid.UUID) *entities.Product {
	for i, product := range h.products {
		if product.ID == productID {
			return &h.products[i]
		}
	}
	return nil
}
