package repositories

import (
	"context"

	"github.com/AndrivA89/orders/internal/domain/repositories"

	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) repositories.TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, repos repositories.TransactionalRepositories) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := repositories.TransactionalRepositories{
			OrderRepository:   NewOrderRepository(tx),
			ProductRepository: NewProductRepository(tx),
			UserRepository:    NewUserRepository(tx),
		}

		return fn(ctx, repos)
	})
}
