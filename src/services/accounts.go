package services

import (
	"account-summary/src/internal/libs"
	"account-summary/src/internal/repositories"
	"account-summary/src/internal/services"
	"account-summary/src/models"
	"context"
	"errors"
	"log"
	"time"
)

type accountsService struct {
	idGenerator       libs.IdGeneratorInterface
	emailSender       libs.EmailSenderInterface
	accountRepository repositories.AccountRepositoryInterface
}

func NewAccountsService(
	accountRepository repositories.AccountRepositoryInterface,
	idGenerator libs.IdGeneratorInterface,
	emailSender libs.EmailSenderInterface,
) services.AccountsServiceInterface {
	return &accountsService{
		accountRepository: accountRepository,
		idGenerator:       idGenerator,
		emailSender:       emailSender,
	}
}

func (s *accountsService) GetOrCreateAccount(ctx context.Context, accountEmail string) (*models.Account, error) {
	account, err := s.accountRepository.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return nil, err
	}

	if account != nil {
		return account, nil
	}

	account = &models.Account{
		ID:        s.idGenerator.Generate("ACNT"),
		Email:     accountEmail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.accountRepository.UpsertAccount(ctx, *account)
	if err != nil {
		log.Default().Printf("error to upsert account: %v\n", err)
		return nil, err
	}

	return account, nil
}

func (s *accountsService) GetAccountByEmail(ctx context.Context, accountEmail string) (*models.Account, error) {
	account, err := s.accountRepository.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return nil, err
	}

	if account == nil {
		log.Default().Printf("account not found: %v\n", accountEmail)
		return nil, errors.New("account not found")
	}

	return account, nil
}

func (s *accountsService) SendSummaryEmail(ctx context.Context, accountEmail string) error {
	account, err := s.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		return err
	}

	err = s.emailSender.SendAccountSummaryEmail(accountEmail, "Transaction Summary", account.Summary)
	if err != nil {
		log.Default().Printf("error to send email: %v\n", err)
		return err
	}

	return nil
}

func (s *accountsService) UpdateAccountSummary(ctx context.Context, account *models.Account, summary models.AccountSummary) error {
	account.Summary = summary
	account.UpdatedAt = time.Now()

	err := s.accountRepository.UpsertAccount(ctx, *account)
	if err != nil {
		log.Default().Printf("error to upsert account: %v\n", err)
		return err
	}

	return nil
}
