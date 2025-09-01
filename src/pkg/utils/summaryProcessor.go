package utils

import (
	"account-summary/src/internal/libs"
	"account-summary/src/models"
)

type summaryProcessor struct {
}

func NewSummaryProcessor() libs.SummaryProcessorInterface {
	return &summaryProcessor{}
}

func (s *summaryProcessor) ProcessSummary(transactions []models.Transaction) (models.AccountSummary, error) {
	summary := models.AccountSummary{}
	summary.Overall.TransactionsPerMonth = make(models.TransactionsPerMonth)
	summary.Debits.TransactionsPerMonth = make(models.TransactionsPerMonth)
	summary.Credits.TransactionsPerMonth = make(models.TransactionsPerMonth)

	for _, t := range transactions {
		if t.Amount > 0 {
			summary.Credits.AddTransaction(t.Amount, t.Date.Month())
		} else {
			summary.Debits.AddTransaction(t.Amount, t.Date.Month())
		}

		summary.Overall.AddTransaction(t.Amount, t.Date.Month())
	}

	summary.Overall.CalculateAverageAmount()
	summary.Credits.CalculateAverageAmount()
	summary.Debits.CalculateAverageAmount()

	return summary, nil
}
