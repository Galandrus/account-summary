package models

import "fmt"

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Transaction string `json:"transaction"`
	Name        string `json:"name"`
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %s, Date: %s, Transaction: %s, Name: %s", t.ID, t.Date, t.Transaction, t.Name)
}

type Transactions []Transaction
