package main

import (
	"log"

	appServices "github.com/AndrivA89/orders/internal/application/services"
	"github.com/AndrivA89/orders/internal/infrastructure/config"
	"github.com/AndrivA89/orders/internal/infrastructure/database"
	infraRepos "github.com/AndrivA89/orders/internal/infrastructure/repositories"
	"github.com/AndrivA89/orders/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Выполняем миграции (в продакшене я бы использовал отдельные миграции, конечно)
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	r := gin.Default()

	userRepo := infraRepos.NewUserRepository(db.DB)
	productRepo := infraRepos.NewProductRepository(db.DB)
	orderRepo := infraRepos.NewOrderRepository(db.DB)

	userService := appServices.NewUserService(userRepo)
	productService := appServices.NewProductService(productRepo)
	orderService := appServices.NewOrderService(orderRepo, userRepo, productRepo)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

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
