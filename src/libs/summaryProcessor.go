package libs

import "account-summary/src/models"

type SummaryProcessor interface {
	ProcessSummary(transactions []models.Transaction) (models.TransactionSummary, error)
}

type summaryProcessor struct {
}

func NewSummaryProcessor() SummaryProcessor {
	return &summaryProcessor{}
}

func (s *summaryProcessor) ProcessSummary(transactions []models.Transaction) (models.TransactionSummary, error) {
	return models.TransactionSummary{}, nil
}
