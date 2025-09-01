package repositories

import (
	"account-summary/src/models"
	"context"
)

type TransactionRepositoryInterface interface {
	GetTransactions(ctx context.Context, accountId string) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction) error
}
