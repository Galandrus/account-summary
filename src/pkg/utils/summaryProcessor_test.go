package utils

import (
	"fmt"
	"testing"
	"time"

	"account-summary/src/models"
)

func TestNewSummaryProcessor(t *testing.T) {
	processor := NewSummaryProcessor()
	if processor == nil {
		t.Error("NewSummaryProcessor should not return nil")
	}
}

func TestProcessSummary_EmptyTransactions(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if summary.Overall.TotalAmount != 0 {
		t.Errorf("Expected overall total amount to be 0, got %f", summary.Overall.TotalAmount)
	}

	if summary.Overall.TotalTransactions != 0 {
		t.Errorf("Expected overall total transactions to be 0, got %d", summary.Overall.TotalTransactions)
	}

	if summary.Overall.AverageAmount != 0 {
		t.Errorf("Expected overall average amount to be 0, got %f", summary.Overall.AverageAmount)
	}

	if summary.Credits.TotalAmount != 0 {
		t.Errorf("Expected credits total amount to be 0, got %f", summary.Credits.TotalAmount)
	}

	if summary.Debits.TotalAmount != 0 {
		t.Errorf("Expected debits total amount to be 0, got %f", summary.Debits.TotalAmount)
	}
}

func TestProcessSummary_OnlyCredits(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{
		{
			ID:     "1",
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount: 100.0,
			Name:   "Credit 1",
		},
		{
			ID:     "2",
			Date:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount: 200.0,
			Name:   "Credit 2",
		},
		{
			ID:     "3",
			Date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Amount: 300.0,
			Name:   "Credit 3",
		},
	}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verificar overall
	expectedOverallTotal := 600.0
	if summary.Overall.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected overall total amount to be %f, got %f", expectedOverallTotal, summary.Overall.TotalAmount)
	}

	if summary.Overall.TotalTransactions != 3 {
		t.Errorf("Expected overall total transactions to be 3, got %d", summary.Overall.TotalTransactions)
	}

	expectedOverallAverage := 200.0
	if summary.Overall.AverageAmount != expectedOverallAverage {
		t.Errorf("Expected overall average amount to be %f, got %f", expectedOverallAverage, summary.Overall.AverageAmount)
	}

	// Verificar credits
	if summary.Credits.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected credits total amount to be %f, got %f", expectedOverallTotal, summary.Credits.TotalAmount)
	}

	if summary.Credits.TotalTransactions != 3 {
		t.Errorf("Expected credits total transactions to be 3, got %d", summary.Credits.TotalTransactions)
	}

	if summary.Credits.AverageAmount != expectedOverallAverage {
		t.Errorf("Expected credits average amount to be %f, got %f", expectedOverallAverage, summary.Credits.AverageAmount)
	}

	// Verificar debits
	if summary.Debits.TotalAmount != 0 {
		t.Errorf("Expected debits total amount to be 0, got %f", summary.Debits.TotalAmount)
	}

	if summary.Debits.TotalTransactions != 0 {
		t.Errorf("Expected debits total transactions to be 0, got %d", summary.Debits.TotalTransactions)
	}

	// Verificar transacciones por mes
	if len(summary.Credits.TransactionsPerMonth) != 2 {
		t.Errorf("Expected 2 months in credits, got %d", len(summary.Credits.TransactionsPerMonth))
	}

	januaryStats := summary.Credits.TransactionsPerMonth[time.January]
	if januaryStats.TotalAmount != 300.0 {
		t.Errorf("Expected January credits total to be 300.0, got %f", januaryStats.TotalAmount)
	}

	if januaryStats.TotalTransactions != 2 {
		t.Errorf("Expected January credits transactions to be 2, got %d", januaryStats.TotalTransactions)
	}

	februaryStats := summary.Credits.TransactionsPerMonth[time.February]
	if februaryStats.TotalAmount != 300.0 {
		t.Errorf("Expected February credits total to be 300.0, got %f", februaryStats.TotalAmount)
	}

	if februaryStats.TotalTransactions != 1 {
		t.Errorf("Expected February credits transactions to be 1, got %d", februaryStats.TotalTransactions)
	}
}

func TestProcessSummary_OnlyDebits(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{
		{
			ID:     "1",
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount: -100.0,
			Name:   "Debit 1",
		},
		{
			ID:     "2",
			Date:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount: -200.0,
			Name:   "Debit 2",
		},
		{
			ID:     "3",
			Date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Amount: -300.0,
			Name:   "Debit 3",
		},
	}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verificar overall
	expectedOverallTotal := -600.0
	if summary.Overall.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected overall total amount to be %f, got %f", expectedOverallTotal, summary.Overall.TotalAmount)
	}

	if summary.Overall.TotalTransactions != 3 {
		t.Errorf("Expected overall total transactions to be 3, got %d", summary.Overall.TotalTransactions)
	}

	expectedOverallAverage := -200.0
	if summary.Overall.AverageAmount != expectedOverallAverage {
		t.Errorf("Expected overall average amount to be %f, got %f", expectedOverallAverage, summary.Overall.AverageAmount)
	}

	// Verificar debits
	if summary.Debits.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected debits total amount to be %f, got %f", expectedOverallTotal, summary.Debits.TotalAmount)
	}

	if summary.Debits.TotalTransactions != 3 {
		t.Errorf("Expected debits total transactions to be 3, got %d", summary.Debits.TotalTransactions)
	}

	if summary.Debits.AverageAmount != expectedOverallAverage {
		t.Errorf("Expected debits average amount to be %f, got %f", expectedOverallAverage, summary.Debits.AverageAmount)
	}

	// Verificar credits
	if summary.Credits.TotalAmount != 0 {
		t.Errorf("Expected credits total amount to be 0, got %f", summary.Credits.TotalAmount)
	}

	if summary.Credits.TotalTransactions != 0 {
		t.Errorf("Expected credits total transactions to be 0, got %d", summary.Credits.TotalTransactions)
	}
}

func TestProcessSummary_MixedTransactions(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{
		{
			ID:     "1",
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount: 100.0,
			Name:   "Credit 1",
		},
		{
			ID:     "2",
			Date:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Amount: -50.0,
			Name:   "Debit 1",
		},
		{
			ID:     "3",
			Date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Amount: 200.0,
			Name:   "Credit 2",
		},
		{
			ID:     "4",
			Date:   time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			Amount: -75.0,
			Name:   "Debit 2",
		},
	}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verificar overall
	expectedOverallTotal := 175.0 // 100 + (-50) + 200 + (-75)
	if summary.Overall.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected overall total amount to be %f, got %f", expectedOverallTotal, summary.Overall.TotalAmount)
	}

	if summary.Overall.TotalTransactions != 4 {
		t.Errorf("Expected overall total transactions to be 4, got %d", summary.Overall.TotalTransactions)
	}

	expectedOverallAverage := 43.75 // 175 / 4
	if summary.Overall.AverageAmount != expectedOverallAverage {
		t.Errorf("Expected overall average amount to be %f, got %f", expectedOverallAverage, summary.Overall.AverageAmount)
	}

	// Verificar credits
	expectedCreditsTotal := 300.0 // 100 + 200
	if summary.Credits.TotalAmount != expectedCreditsTotal {
		t.Errorf("Expected credits total amount to be %f, got %f", expectedCreditsTotal, summary.Credits.TotalAmount)
	}

	if summary.Credits.TotalTransactions != 2 {
		t.Errorf("Expected credits total transactions to be 2, got %d", summary.Credits.TotalTransactions)
	}

	expectedCreditsAverage := 150.0 // 300 / 2
	if summary.Credits.AverageAmount != expectedCreditsAverage {
		t.Errorf("Expected credits average amount to be %f, got %f", expectedCreditsAverage, summary.Credits.AverageAmount)
	}

	// Verificar debits
	expectedDebitsTotal := -125.0 // -50 + (-75)
	if summary.Debits.TotalAmount != expectedDebitsTotal {
		t.Errorf("Expected debits total amount to be %f, got %f", expectedDebitsTotal, summary.Debits.TotalAmount)
	}

	if summary.Debits.TotalTransactions != 2 {
		t.Errorf("Expected debits total transactions to be 2, got %d", summary.Debits.TotalTransactions)
	}

	expectedDebitsAverage := -62.5 // -125 / 2
	if summary.Debits.AverageAmount != expectedDebitsAverage {
		t.Errorf("Expected debits average amount to be %f, got %f", expectedDebitsAverage, summary.Debits.AverageAmount)
	}
}

func TestProcessSummary_ZeroAmountTransactions(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{
		{
			ID:     "1",
			Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Amount: 0.0,
			Name:   "Zero Transaction",
		},
	}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Las transacciones con monto 0 se consideran como débitos (amount > 0 es falso)
	if summary.Overall.TotalAmount != 0 {
		t.Errorf("Expected overall total amount to be 0, got %f", summary.Overall.TotalAmount)
	}

	if summary.Overall.TotalTransactions != 1 {
		t.Errorf("Expected overall total transactions to be 1, got %d", summary.Overall.TotalTransactions)
	}

	if summary.Credits.TotalAmount != 0 {
		t.Errorf("Expected credits total amount to be 0, got %f", summary.Credits.TotalAmount)
	}

	if summary.Credits.TotalTransactions != 0 {
		t.Errorf("Expected credits total transactions to be 0, got %d", summary.Credits.TotalTransactions)
	}

	if summary.Debits.TotalAmount != 0 {
		t.Errorf("Expected debits total amount to be 0, got %f", summary.Debits.TotalAmount)
	}

	if summary.Debits.TotalTransactions != 1 {
		t.Errorf("Expected debits total transactions to be 1, got %d", summary.Debits.TotalTransactions)
	}
}

func TestProcessSummary_LargeDataset(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := make([]models.Transaction, 1000)

	// Crear transacciones alternando entre positivas y negativas
	for i := 0; i < 1000; i++ {
		amount := float64(i + 1)
		if i%2 == 0 {
			amount = -amount // Hacer negativos los pares (índices 0, 2, 4, ...)
		}
		transactions[i] = models.Transaction{
			ID:     string(rune(i)),
			Date:   time.Date(2024, time.Month((i%12)+1), (i%28)+1, 0, 0, 0, 0, time.UTC),
			Amount: amount,
			Name:   fmt.Sprintf("Transaction %d", i),
		}
	}

	summary, err := processor.ProcessSummary(transactions)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if summary.Overall.TotalTransactions != 1000 {
		t.Errorf("Expected overall total transactions to be 1000, got %d", summary.Overall.TotalTransactions)
	}

	if summary.Credits.TotalTransactions != 500 {
		t.Errorf("Expected credits total transactions to be 500, got %d", summary.Credits.TotalTransactions)
	}

	if summary.Debits.TotalTransactions != 500 {
		t.Errorf("Expected debits total transactions to be 500, got %d", summary.Debits.TotalTransactions)
	}

	// Calcular los montos esperados manualmente
	// Índices impares (1, 3, 5, ..., 999): 1+3+5+...+999 = 250500
	// Índices pares (0, 2, 4, ..., 998): -1-3-5-...-999 = -250000
	expectedCreditsTotal := float64(250500)
	expectedDebitsTotal := float64(-250000)

	if summary.Credits.TotalAmount != expectedCreditsTotal {
		t.Errorf("Expected credits total amount to be %f, got %f", expectedCreditsTotal, summary.Credits.TotalAmount)
	}

	if summary.Debits.TotalAmount != expectedDebitsTotal {
		t.Errorf("Expected debits total amount to be %f, got %f", expectedDebitsTotal, summary.Debits.TotalAmount)
	}

	expectedOverallTotal := expectedCreditsTotal + expectedDebitsTotal
	if summary.Overall.TotalAmount != expectedOverallTotal {
		t.Errorf("Expected overall total amount to be %f, got %f", expectedOverallTotal, summary.Overall.TotalAmount)
	}
}
