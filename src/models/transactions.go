package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        string
	Date      time.Time
	Amount    float64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %s, Date: %s, Amount: %0.2f, Name: %s", t.ID, t.Date, t.Amount, t.Name)
}

type Transactions []Transaction
