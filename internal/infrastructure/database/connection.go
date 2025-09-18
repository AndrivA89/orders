package database

import (
	"github.com/AndrivA89/orders/internal/infrastructure/config"
	"github.com/AndrivA89/orders/internal/infrastructure/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	DB *gorm.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &Connection{DB: db}, nil
}

func (c *Connection) AutoMigrate() error {
	return c.DB.AutoMigrate(
		&models.UserModel{},
		&models.ProductModel{},
		&models.OrderModel{},
		&models.OrderItemModel{},
	)
}

func (c *Connection) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
