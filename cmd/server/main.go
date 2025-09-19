package main

import (
	"github.com/sirupsen/logrus"

	appServices "github.com/AndrivA89/orders/internal/application/services"
	"github.com/AndrivA89/orders/internal/infrastructure/config"
	"github.com/AndrivA89/orders/internal/infrastructure/database"
	infraRepos "github.com/AndrivA89/orders/internal/infrastructure/repositories"
	"github.com/AndrivA89/orders/internal/transport/http/handlers"
	"github.com/AndrivA89/orders/internal/transport/http/router"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	if cfg.Logger.Level == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Initialize database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize repositories
	userRepo := infraRepos.NewUserRepository(db.DB)
	productRepo := infraRepos.NewProductRepository(db.DB)
	orderRepo := infraRepos.NewOrderRepository(db.DB)
	txManager := infraRepos.NewTransactionManager(db.DB)

	// Initialize services
	userService := appServices.NewUserService(userRepo)
	productService := appServices.NewProductService(productRepo)
	orderService := appServices.NewOrderService(orderRepo, userRepo, productRepo, txManager)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router with middleware
	r := router.NewRouter(userHandler, productHandler, orderHandler, logger)
	engine := r.SetupRoutes()

	logger.Infof("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
