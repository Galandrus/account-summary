package services

import (
	"account-summary/src/libs"
	"account-summary/src/models"
	"account-summary/src/repository"
	"context"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"
)

type TransactionsService interface {
	GetTransactions(ctx context.Context, accountId string) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction, accountId string) error
	LoadTransactions(ctx context.Context, path string, accountId string) error
	GetSummary(ctx context.Context, accountId string) (models.TransactionSummary, error)
	SendEmail(ctx context.Context, accountId string) ([]byte, error)
}

type transactionsService struct {
	csvReader             libs.CsvReader
	summaryProcessor      libs.SummaryProcessor
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
}

func NewTransactionsService(
	transactionRepository repository.TransactionRepository,
	accountRepository repository.AccountRepository,
	csvReader libs.CsvReader,
	summaryProcessor libs.SummaryProcessor,
) TransactionsService {
	return &transactionsService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		csvReader:             csvReader,
		summaryProcessor:      summaryProcessor,
	}
}

func (s *transactionsService) GetTransactions(ctx context.Context, accountId string) ([]models.Transaction, error) {
	return s.transactionRepository.GetTransactions(ctx, accountId)
}

func (s *transactionsService) CreateTransactions(ctx context.Context, transactions []models.Transaction, accountId string) error {
	for idx, t := range transactions {
		transaction := t
		transaction.AccountId = accountId
		transaction.CreatedAt = time.Now()
		transaction.UpdatedAt = time.Now()
		transactions[idx] = transaction
	}

	return s.transactionRepository.CreateTransactions(ctx, transactions)
}

func (s *transactionsService) LoadTransactions(ctx context.Context, path string, accountId string) error {
	transactions, err := s.csvReader.LoadTransactions(path)
	if err != nil {
		log.Default().Printf("error to load transactions: %v\n", err)
		return err
	}

	err = s.CreateTransactions(ctx, transactions, accountId)
	if err != nil {
		log.Default().Printf("error to create transactions: %v\n", err)
		return err
	}

	return s.UpdateSummary(ctx, accountId)
}

func (s *transactionsService) UpdateSummary(ctx context.Context, accountId string) error {
	summary, err := s.GetSummary(ctx, accountId)
	if err != nil {
		log.Default().Printf("error to update summary: %v\n", err)
		return err
	}

	return s.UpsertAccount(ctx, accountId, summary)
}

func (s *transactionsService) GetSummary(ctx context.Context, accountId string) (models.TransactionSummary, error) {
	transactions, err := s.GetTransactions(ctx, accountId)
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

func (s *transactionsService) UpsertAccount(ctx context.Context, accountId string, summary models.TransactionSummary) error {
	account, err := s.accountRepository.GetAccountById(ctx, accountId)
	if err != nil {
		log.Default().Printf("error to get account: %v\n", err)
		return err
	}

	if account == nil {
		account = &models.Account{
			ID:        accountId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	account.Summary = summary

	err = s.accountRepository.UpsertAccount(ctx, *account)
	if err != nil {
		log.Default().Printf("error to upsert account: %v\n", err)
		return err
	}

	return nil
}

func (s *transactionsService) SendEmail(ctx context.Context, accountId string) ([]byte, error) {
	summary, err := s.GetSummary(ctx, accountId)
	if err != nil {
		return nil, err
	}

	// Crear estructura de datos para el template
	templateData := struct {
		models.TransactionSummary
		GeneratedDate string
	}{
		TransactionSummary: summary,
		GeneratedDate:      time.Now().Format("January 2, 2006 at 3:04 PM"),
	}

	// Parsear el template
	tmpl, err := template.ParseFiles("src/templates/summary.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	// Ejecutar el template y capturar el resultado
	var htmlBuffer strings.Builder
	err = tmpl.Execute(&htmlBuffer, templateData)
	if err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	return []byte(htmlBuffer.String()), nil
}
