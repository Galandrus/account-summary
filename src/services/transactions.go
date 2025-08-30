package services

import (
	"account-summary/src/libs"
	"account-summary/src/models"
	"account-summary/src/repository"
	"context"
)

type TransactionsService interface {
	GetTransactions(ctx context.Context) ([]models.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []models.Transaction) error
	LoadTransactions(ctx context.Context, path string) error
	GetSummary(ctx context.Context) (models.TransactionSummary, error)
}

type transactionsService struct {
	csvReader        libs.CsvReader
	summaryProcessor libs.SummaryProcessor
	repository       repository.TransactionRepository
}

func NewTransactionsService(repository repository.TransactionRepository, csvReader libs.CsvReader, summaryProcessor libs.SummaryProcessor) TransactionsService {
	return &transactionsService{repository: repository, csvReader: csvReader, summaryProcessor: summaryProcessor}
}

func (s *transactionsService) GetTransactions(ctx context.Context) ([]models.Transaction, error) {
	return s.repository.GetTransactions(ctx)
}

func (s *transactionsService) CreateTransactions(ctx context.Context, transactions []models.Transaction) error {
	return s.repository.CreateTransactions(ctx, transactions)
}

func (s *transactionsService) LoadTransactions(ctx context.Context, path string) error {
	transactions, err := s.csvReader.LoadTransactions(path)
	if err != nil {
		return err
	}

	return s.repository.CreateTransactions(ctx, transactions)
}

func (s *transactionsService) GetSummary(ctx context.Context) (models.TransactionSummary, error) {
	transactions, err := s.GetTransactions(ctx)
	if err != nil {
		return models.TransactionSummary{}, err
	}

	summary, err := s.summaryProcessor.ProcessSummary(transactions)
	if err != nil {
		return models.TransactionSummary{}, err
	}

	return summary, nil
}
