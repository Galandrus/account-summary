package models

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type AccountSummary struct {
	Overall SummaryStats `json:"overall" bson:"overall"`
	Debits  SummaryStats `json:"debits" bson:"debits"`
	Credits SummaryStats `json:"credits" bson:"credits"`
}

func (s *AccountSummary) String() string {
	return fmt.Sprintf("Overall: \n%s\nDebits: \n%s\nCredits: \n%s", s.Overall.String(), s.Debits.String(), s.Credits.String())
}

type SummaryStats struct {
	TotalAmount          float64              `json:"totalAmount" bson:"totalAmount"`
	TotalTransactions    int                  `json:"totalTransactions" bson:"totalTransactions"`
	AverageAmount        float64              `json:"averageAmount" bson:"averageAmount"`
	TransactionsPerMonth TransactionsPerMonth `json:"transactionsPerMonth" bson:"transactionsPerMonth"`
}

func (s *SummaryStats) CalculateAverageAmount() {
	if s.TotalTransactions == 0 {
		s.AverageAmount = 0
		return
	}

	s.AverageAmount = s.TotalAmount / float64(s.TotalTransactions)

	for month, stats := range s.TransactionsPerMonth {
		stats.CalculateAverageAmount()
		s.TransactionsPerMonth[month] = stats
	}
}

func (s *SummaryStats) AddTransaction(amount float64, month time.Month) {
	s.TotalAmount += amount
	s.TotalTransactions++

	stats, ok := s.TransactionsPerMonth[month]
	if !ok {
		stats = SummaryStats{}
	}

	stats.TotalAmount += amount
	stats.TotalTransactions++

	s.TransactionsPerMonth[month] = stats
}

func (s *SummaryStats) String() string {
	return fmt.Sprintf("  Total Amount: $%0.2f\n  Total Transactions: %d\n  Average Amount: $%0.2f\n  Transactions Per Month: \n%s", s.TotalAmount, s.TotalTransactions, s.AverageAmount, s.TransactionsPerMonth.String())
}

type TransactionsPerMonth map[time.Month]SummaryStats

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
		months = append(months, fmt.Sprintf("    Transactions in %s: %d (Average: $%0.2f)", month.String(), count.TotalTransactions, count.AverageAmount))
	}

	return strings.Join(months, "\n")
}
