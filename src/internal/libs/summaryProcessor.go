package libs

import "account-summary/src/models"

type SummaryProcessorInterface interface {
	ProcessSummary(transactions []models.Transaction) (models.AccountSummary, error)
}
