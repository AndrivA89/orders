package repositories

//go:generate mockgen -source=transaction_manager.go -destination=mocks/transaction_manager_mock.go -package=mocks

import (
	"context"
)

type TransactionManager interface {
	// WithTransaction выполняет функцию в рамках транзакции
	WithTransaction(ctx context.Context, fn func(ctx context.Context, repos TransactionalRepositories) error) error
}

type TransactionalRepositories struct {
	OrderRepository   OrderRepository
	ProductRepository ProductRepository
	UserRepository    UserRepository
}
