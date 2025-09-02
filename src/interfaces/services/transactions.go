package services

import (
	"account-summary/src/models"
	"context"
)

type TransactionsServiceInterface interface {
	GetTransactionsByEmail(ctx context.Context, accountEmail string) ([]models.Transaction, error)
	GetTransactionsByID(ctx context.Context, accountID string) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction, account *models.Account) error
	LoadTransactions(ctx context.Context, path string, accountEmail string) error
}

type AccountsServiceInterface interface {
	SendSummaryEmail(ctx context.Context, accountEmail string) error
	GetAccountByEmail(ctx context.Context, accountEmail string) (*models.Account, error)
	GetOrCreateAccount(ctx context.Context, accountEmail string) (*models.Account, error)
	UpdateAccountSummary(ctx context.Context, account *models.Account, summary models.AccountSummary) error
}

type SummaryServiceInterface interface {
	GenerateSummary(ctx context.Context, account *models.Account) (models.AccountSummary, error)
	UpdateSummary(ctx context.Context, account *models.Account) error
}
