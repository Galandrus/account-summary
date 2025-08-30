package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        string    `json:"id" bson:"id"`
	Date      time.Time `json:"date" bson:"date"`
	Amount    float64   `json:"amount" bson:"amount"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	AccountId string    `json:"accountId" bson:"accountId"`
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %s, Date: %s, Amount: %0.2f, Name: %s", t.ID, t.Date, t.Amount, t.Name)
}

type Transactions []Transaction
