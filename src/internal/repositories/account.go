package repositories

import (
	"account-summary/src/models"
	"context"
)

type AccountRepositoryInterface interface {
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountById(ctx context.Context, id string) (*models.Account, error)
	UpsertAccount(ctx context.Context, account models.Account) error
}
