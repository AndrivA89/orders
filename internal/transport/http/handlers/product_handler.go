package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndrivA89/orders/internal/domain/services"
	"github.com/AndrivA89/orders/internal/transport/http/dto"
	"github.com/AndrivA89/orders/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), req.ToServiceRequest())
	if err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToProductResponse(product))
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	products, err := h.productService.GetProducts(c.Request.Context(), limit, offset)
	if err != nil {
		middleware.HandleInternalError(c, err)
		return
	}

	responses := make([]*dto.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = dto.ToProductResponse(product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": responses,
		"total":    len(products),
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		middleware.HandleValidationError(c, err)
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		middleware.HandleNotFoundError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToProductResponse(product))
}
