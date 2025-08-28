package models

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type TransactionSummary struct {
	Overall SummaryStats
	Debits  SummaryStats
	Credits SummaryStats
}

func (s *TransactionSummary) String() string {
	return fmt.Sprintf("Overall: \n%s\nDebits: \n%s\nCredits: \n%s", s.Overall.String(), s.Debits.String(), s.Credits.String())
}

type SummaryStats struct {
	TotalAmount          float64
	TotalTransactions    int
	AverageAmount        float64
	TransactionsPerMonth TransactionsPerMonth
}

func (s *SummaryStats) CalculateAverageAmount() {
	s.AverageAmount = s.TotalAmount / float64(s.TotalTransactions)
}

func (s *SummaryStats) String() string {
	return fmt.Sprintf("  Total Amount: %0.2f\n  Total Transactions: %d\n  Average Amount: %0.2f\n  Transactions Per Month: \n%s", s.TotalAmount, s.TotalTransactions, s.AverageAmount, s.TransactionsPerMonth.String())
}

type TransactionsPerMonth map[time.Month]int

func (t *TransactionsPerMonth) Sort() []time.Month {
	months := []time.Month{}

	for month := range *t {
		months = append(months, month)
	}

	slices.Sort(months)

	return months
}

func (t *TransactionsPerMonth) String() string {
	months := []string{}

	for _, month := range t.Sort() {
		count := (*t)[month]
		months = append(months, fmt.Sprintf("    Transactions in %s: %d", month.String(), count))
	}

	return strings.Join(months, "\n")
}
