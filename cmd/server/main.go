package main

import (
	"github.com/AndrivA89/orders/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userHandler := handlers.NewUserHandler()
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler(productHandler, userHandler)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "orders",
		})
	})

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/:user_id/orders", orderHandler.GetOrdersByUser)
		}

		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
		}

		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PATCH("/:id/confirm", orderHandler.ConfirmOrder)
			orders.PATCH("/:id/cancel", orderHandler.CancelOrder)
		}
	}

	r.Run(":8080")
}
