package models

import "fmt"

type TransactionSummary struct {
	Overall SummaryStats
	Debits  SummaryStats
	Credits SummaryStats
}

type SummaryStats struct {
	TotalAmount          float64
	TotalTransactions    int
	AverageAmount        float64
	TransactionsPerMonth map[string]int
}

func (s *TransactionSummary) String() string {
	return fmt.Sprintf("Overall: %v, Debits: %v, Credits: %v", s.Overall, s.Debits, s.Credits)
}
