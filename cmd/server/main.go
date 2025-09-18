package main

import (
	"github.com/AndrivA89/orders/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userHandler := handlers.NewUserHandler()
	productHandler := handlers.NewProductHandler()

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
		}

		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
		}

		v1.GET("/orders", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Orders endpoint - coming soon",
			})
		})
	}

	r.Run(":8080")
}
