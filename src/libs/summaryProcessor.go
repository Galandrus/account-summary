package libs

import (
	"account-summary/src/models"
)

type SummaryProcessor interface {
	ProcessSummary(transactions []models.Transaction) (models.TransactionSummary, error)
}

type summaryProcessor struct {
}

func NewSummaryProcessor() SummaryProcessor {
	return &summaryProcessor{}
}

func (s *summaryProcessor) ProcessSummary(transactions []models.Transaction) (models.TransactionSummary, error) {
	summary := models.TransactionSummary{}
	summary.Overall.TransactionsPerMonth = make(models.TransactionsPerMonth)
	summary.Debits.TransactionsPerMonth = make(models.TransactionsPerMonth)
	summary.Credits.TransactionsPerMonth = make(models.TransactionsPerMonth)

	for _, t := range transactions {
		if t.Amount > 0 {
			summary.Credits.TotalAmount += t.Amount
			summary.Credits.TotalTransactions++
		} else {
			summary.Debits.TotalAmount += t.Amount
			summary.Debits.TotalTransactions++
		}

		summary.Overall.TransactionsPerMonth[t.Date.Month()]++
		summary.Debits.TransactionsPerMonth[t.Date.Month()]++
		summary.Credits.TransactionsPerMonth[t.Date.Month()]++
	}

	summary.Overall.TotalAmount = summary.Credits.TotalAmount + summary.Debits.TotalAmount
	summary.Overall.TotalTransactions = summary.Credits.TotalTransactions + summary.Debits.TotalTransactions
	summary.Overall.CalculateAverageAmount()
	summary.Credits.CalculateAverageAmount()
	summary.Debits.CalculateAverageAmount()

	return summary, nil
}
