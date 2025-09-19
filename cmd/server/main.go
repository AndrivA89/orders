package main

import (
	"fmt"
	"log"

	"github.com/AndrivA89/orders/internal/application/services"
	"github.com/AndrivA89/orders/internal/infrastructure/config"
	"github.com/AndrivA89/orders/internal/infrastructure/database"
	"github.com/AndrivA89/orders/internal/infrastructure/repositories"
	"github.com/AndrivA89/orders/internal/infrastructure/telemetry"
	"github.com/AndrivA89/orders/internal/transport/http/handlers"
	"github.com/AndrivA89/orders/internal/transport/http/router"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	logger := setupLogger(cfg)

	cleanup, err := telemetry.InitTracing("orders-service")
	if err != nil {
		logger.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer cleanup()

	dbConn, err := database.NewConnection(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		if err = dbConn.Close(); err != nil {
			logger.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// Выполняем миграции (в продакшене я бы использовал отдельные миграции, конечно)
	if err = dbConn.AutoMigrate(); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := repositories.NewUserRepository(dbConn.DB)
	productRepo := repositories.NewProductRepository(dbConn.DB)
	orderRepo := repositories.NewOrderRepository(dbConn.DB)

	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	txManager := repositories.NewTransactionManager(dbConn.DB)
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo, txManager)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	appRouter := router.NewRouter(userHandler, productHandler, orderHandler, logger)
	ginRouter := appRouter.SetupRoutes()

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Starting server on %s", address)

	if err = ginRouter.Run(address); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		log.Printf("Invalid log level %s, using INFO", cfg.Logger.Level)
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}
