package libs

import "account-summary/src/models"

type CsvReaderInterface interface {
	LoadTransactions(path string) ([]models.Transaction, error)
}
