package services

import (
	"account-summary/src/libs"
	"account-summary/src/models"
	"account-summary/src/repository"
	"context"
	"errors"
	"log"
	"time"
)

type TransactionsServiceInterface interface {
	GetTransactionsByEmail(ctx context.Context, accountEmail string) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction, account *models.Account) error
	LoadTransactions(ctx context.Context, path string, accountEmail string) error
	GetSummary(ctx context.Context, account *models.Account) (models.TransactionSummary, error)
	SendEmail(ctx context.Context, accountEmail string) error
	GetOrCreateAccount(ctx context.Context, accountEmail string) (*models.Account, error)
	UpdateAccountSummary(ctx context.Context, account *models.Account, summary models.TransactionSummary) error
	UpdateSummary(ctx context.Context, account *models.Account) error
	GetSummaryByEmail(ctx context.Context, accountEmail string) (models.TransactionSummary, error)
}

type transactionsService struct {
	csvReader             libs.CsvReader
	summaryProcessor      libs.SummaryProcessor
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
	emailSender           libs.EmailSender
	idGenerator           libs.IdGenerator
}

func NewTransactionsService(
	transactionRepository repository.TransactionRepository,
	accountRepository repository.AccountRepository,
	csvReader libs.CsvReader,
	summaryProcessor libs.SummaryProcessor,
	emailSender libs.EmailSender,
	idGenerator libs.IdGenerator,
) TransactionsServiceInterface {
	return &transactionsService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		csvReader:             csvReader,
		summaryProcessor:      summaryProcessor,
		emailSender:           emailSender,
		idGenerator:           idGenerator,
	}
}

func (s *transactionsService) GetTransactionsByEmail(ctx context.Context, accountEmail string) ([]models.Transaction, error) {
	account, err := s.accountRepository.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return nil, err
	}

	if account == nil {
		log.Default().Printf("account not found: %v\n", accountEmail)
		return nil, errors.New("account not found")
	}

	return s.transactionRepository.GetTransactions(ctx, account.ID)
}

func (s *transactionsService) CreateTransactions(ctx context.Context, transactions []models.Transaction, account *models.Account) error {
	for idx, t := range transactions {
		transaction := t
		transaction.AccountId = account.ID
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = time.Now()
		transactions[idx] = transaction
	}

	return s.transactionRepository.CreateTransactions(ctx, transactions)
}

func (s *transactionsService) LoadTransactions(ctx context.Context, path string, accountEmail string) error {
	transactions, err := s.csvReader.LoadTransactions(path)
	if err != nil {
		log.Default().Printf("error to load transactions: %v\n", err)
		return err
	}

	account, err := s.GetOrCreateAccount(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get or create account: %v\n", err)
		return err
	}

	err = s.CreateTransactions(ctx, transactions, account)
	if err != nil {
		log.Default().Printf("error to create transactions: %v\n", err)
		return err
	}

	return s.UpdateSummary(ctx, account)
}

func (s *transactionsService) UpdateSummary(ctx context.Context, account *models.Account) error {
	summary, err := s.GetSummary(ctx, account)
	if err != nil {
		log.Default().Printf("error to update summary: %v\n", err)
		return err
	}

	return s.UpdateAccountSummary(ctx, account, summary)
}

func (s *transactionsService) GetSummaryByEmail(ctx context.Context, accountEmail string) (models.TransactionSummary, error) {
	account, err := s.accountRepository.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return models.TransactionSummary{}, err
	}

	if account == nil {
		log.Default().Printf("account not found: %v\n", accountEmail)
		return models.TransactionSummary{}, errors.New("account not found")
	}

	return s.GetSummary(ctx, account)
}

func (s *transactionsService) GetSummary(ctx context.Context, account *models.Account) (models.TransactionSummary, error) {
	transactions, err := s.transactionRepository.GetTransactions(ctx, account.ID)
	if err != nil {
		log.Default().Printf("error to get summary: %v\n", err)
		return models.TransactionSummary{}, err
	}

	summary, err := s.summaryProcessor.ProcessSummary(transactions)
	if err != nil {
		log.Default().Printf("error to process summary: %v\n", err)
		return models.TransactionSummary{}, err
	}

	return summary, nil
}

func (s *transactionsService) UpdateAccountSummary(ctx context.Context, account *models.Account, summary models.TransactionSummary) error {
	account.Summary = summary
	account.UpdatedAt = time.Now()

	err := s.accountRepository.UpsertAccount(ctx, *account)
	if err != nil {
		log.Default().Printf("error to upsert account: %v\n", err)
		return err
	}

	return nil
}

func (s *transactionsService) SendEmail(ctx context.Context, accountEmail string) error {
	account, err := s.accountRepository.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return err
	}

	err = s.emailSender.SendEmail(accountEmail, "Transaction Summary", account.Summary)
	if err != nil {
		log.Default().Printf("error to send email: %v\n", err)
		return err
	}

	return nil
}

func (s *transactionsService) GetOrCreateAccount(ctx context.Context, accountEmail string) (*models.Account, error) {
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
