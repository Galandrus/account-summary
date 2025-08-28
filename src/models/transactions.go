package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID     string    `json:"id"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	Name   string    `json:"name"`
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %s, Date: %s, Amount: %0.2f, Name: %s", t.ID, t.Date, t.Amount, t.Name)
}

type Transactions []Transaction
