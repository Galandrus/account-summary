package services

import (
	"account-summary/src/internal/libs"
	"account-summary/src/internal/repositories"
	"account-summary/src/internal/services"
	"account-summary/src/models"
	"context"
	"log"
	"time"
)

type transactionsService struct {
	csvReader             libs.CsvReaderInterface
	transactionRepository repositories.TransactionRepositoryInterface
	accountsService       services.AccountsServiceInterface
	summaryProcessor      libs.SummaryProcessorInterface
}

func NewTransactionsService(
	transactionRepository repositories.TransactionRepositoryInterface,
	accountsService services.AccountsServiceInterface,
	csvReader libs.CsvReaderInterface,
	summaryProcessor libs.SummaryProcessorInterface,
) services.TransactionsServiceInterface {
	return &transactionsService{
		transactionRepository: transactionRepository,
		accountsService:       accountsService,
		csvReader:             csvReader,
		summaryProcessor:      summaryProcessor,
	}
}

func (s *transactionsService) GetTransactionsByEmail(ctx context.Context, accountEmail string) ([]models.Transaction, error) {
	account, err := s.accountsService.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		return nil, err
	}

	return s.GetTransactionsByID(ctx, account.ID)
}

func (s *transactionsService) GetTransactionsByID(ctx context.Context, accountID string) ([]models.Transaction, error) {
	transactions, err := s.transactionRepository.GetTransactions(ctx, accountID)
	if err != nil {
		log.Default().Printf("error to get transactions: %v\n", err)
		return nil, err
	}

	return transactions, nil
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

	account, err := s.accountsService.GetOrCreateAccount(ctx, accountEmail)
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
	summary, err := s.GenerateSummary(ctx, account)
	if err != nil {
		log.Default().Printf("error to update summary: %v\n", err)
		return err
	}

	return s.accountsService.UpdateAccountSummary(ctx, account, summary)
}

func (s *transactionsService) GenerateSummary(ctx context.Context, account *models.Account) (models.AccountSummary, error) {
	transactions, err := s.GetTransactionsByID(ctx, account.ID)
	if err != nil {
		log.Default().Printf("error to get summary: %v\n", err)
		return models.AccountSummary{}, err
	}

	summary, err := s.summaryProcessor.ProcessSummary(transactions)
	if err != nil {
		log.Default().Printf("error to process summary: %v\n", err)
		return models.AccountSummary{}, err
	}

	return summary, nil
}
