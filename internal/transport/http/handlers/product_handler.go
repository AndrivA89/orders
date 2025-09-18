package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AndrivA89/orders/internal/domain/entities"
	"github.com/AndrivA89/orders/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	products []entities.Product
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		products: make([]entities.Product, 0),
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := entities.Product{
		ID:          uuid.New(),
		Description: req.Description,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	h.products = append(h.products, product)

	c.JSON(http.StatusCreated, dto.ToProductResponse(&product))
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	total := len(h.products)
	products := h.products

	if offset >= total {
		products = []entities.Product{}
	} else {
		end := offset + limit
		if end > total {
			end = total
		}
		products = h.products[offset:end]
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = dto.ToProductResponse(&product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": responses,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	for _, product := range h.products {
		if product.ID == productID {
			c.JSON(http.StatusOK, dto.ToProductResponse(&product))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
}
