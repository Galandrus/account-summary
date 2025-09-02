package utils

import (
	"fmt"
	"testing"
	"time"

	"account-summary/src/models"

	"github.com/stretchr/testify/assert"
)

func TestNewSummaryProcessor(t *testing.T) {
	processor := NewSummaryProcessor()
	assert.NotNil(t, processor, "NewSummaryProcessor should not return nil")
}

func TestProcessSummary_EmptyTransactions(t *testing.T) {
	processor := NewSummaryProcessor()
	transactions := []models.Transaction{}

	summary, err := processor.ProcessSummary(transactions)

	assert.NoError(t, err)
	assert.Equal(t, 0.0, summary.Overall.TotalAmount)
	assert.Equal(t, 0, summary.Overall.TotalTransactions)
	assert.Equal(t, 0.0, summary.Overall.AverageAmount)
	assert.Equal(t, 0.0, summary.Credits.TotalAmount)
	assert.Equal(t, 0.0, summary.Debits.TotalAmount)
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

	assert.NoError(t, err)

	// Verificar overall
	expectedOverallTotal := 600.0
	assert.Equal(t, expectedOverallTotal, summary.Overall.TotalAmount)
	assert.Equal(t, 3, summary.Overall.TotalTransactions)

	expectedOverallAverage := 200.0
	assert.Equal(t, expectedOverallAverage, summary.Overall.AverageAmount)

	// Verificar credits
	assert.Equal(t, expectedOverallTotal, summary.Credits.TotalAmount)
	assert.Equal(t, 3, summary.Credits.TotalTransactions)
	assert.Equal(t, expectedOverallAverage, summary.Credits.AverageAmount)

	// Verificar debits
	assert.Equal(t, 0.0, summary.Debits.TotalAmount)
	assert.Equal(t, 0, summary.Debits.TotalTransactions)

	// Verificar transacciones por mes
	assert.Len(t, summary.Credits.TransactionsPerMonth, 2)

	januaryStats := summary.Credits.TransactionsPerMonth[time.January]
	assert.Equal(t, 300.0, januaryStats.TotalAmount)
	assert.Equal(t, 2, januaryStats.TotalTransactions)

	februaryStats := summary.Credits.TransactionsPerMonth[time.February]
	assert.Equal(t, 300.0, februaryStats.TotalAmount)
	assert.Equal(t, 1, februaryStats.TotalTransactions)
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

	assert.NoError(t, err)

	// Verificar overall
	expectedOverallTotal := -600.0
	assert.Equal(t, expectedOverallTotal, summary.Overall.TotalAmount)
	assert.Equal(t, 3, summary.Overall.TotalTransactions)

	expectedOverallAverage := -200.0
	assert.Equal(t, expectedOverallAverage, summary.Overall.AverageAmount)

	// Verificar debits
	assert.Equal(t, expectedOverallTotal, summary.Debits.TotalAmount)
	assert.Equal(t, 3, summary.Debits.TotalTransactions)
	assert.Equal(t, expectedOverallAverage, summary.Debits.AverageAmount)

	// Verificar credits
	assert.Equal(t, 0.0, summary.Credits.TotalAmount)
	assert.Equal(t, 0, summary.Credits.TotalTransactions)
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

	assert.NoError(t, err)

	// Verificar overall
	expectedOverallTotal := 175.0 // 100 + (-50) + 200 + (-75)
	assert.Equal(t, expectedOverallTotal, summary.Overall.TotalAmount)
	assert.Equal(t, 4, summary.Overall.TotalTransactions)

	expectedOverallAverage := 43.75 // 175 / 4
	assert.Equal(t, expectedOverallAverage, summary.Overall.AverageAmount)

	// Verificar credits
	expectedCreditsTotal := 300.0 // 100 + 200
	assert.Equal(t, expectedCreditsTotal, summary.Credits.TotalAmount)
	assert.Equal(t, 2, summary.Credits.TotalTransactions)

	expectedCreditsAverage := 150.0 // 300 / 2
	assert.Equal(t, expectedCreditsAverage, summary.Credits.AverageAmount)

	// Verificar debits
	expectedDebitsTotal := -125.0 // -50 + (-75)
	assert.Equal(t, expectedDebitsTotal, summary.Debits.TotalAmount)
	assert.Equal(t, 2, summary.Debits.TotalTransactions)

	expectedDebitsAverage := -62.5 // -125 / 2
	assert.Equal(t, expectedDebitsAverage, summary.Debits.AverageAmount)
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

	assert.NoError(t, err)

	// Las transacciones con monto 0 se consideran como débitos (amount > 0 es falso)
	assert.Equal(t, 0.0, summary.Overall.TotalAmount)
	assert.Equal(t, 1, summary.Overall.TotalTransactions)
	assert.Equal(t, 0.0, summary.Credits.TotalAmount)
	assert.Equal(t, 0, summary.Credits.TotalTransactions)
	assert.Equal(t, 0.0, summary.Debits.TotalAmount)
	assert.Equal(t, 1, summary.Debits.TotalTransactions)
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

	assert.NoError(t, err)
	assert.Equal(t, 1000, summary.Overall.TotalTransactions)
	assert.Equal(t, 500, summary.Credits.TotalTransactions)
	assert.Equal(t, 500, summary.Debits.TotalTransactions)

	// Calcular los montos esperados manualmente
	// Índices impares (1, 3, 5, ..., 999): 1+3+5+...+999 = 250500
	// Índices pares (0, 2, 4, ..., 998): -1-3-5-...-999 = -250000
	expectedCreditsTotal := float64(250500)
	expectedDebitsTotal := float64(-250000)

	assert.Equal(t, expectedCreditsTotal, summary.Credits.TotalAmount)
	assert.Equal(t, expectedDebitsTotal, summary.Debits.TotalAmount)

	expectedOverallTotal := expectedCreditsTotal + expectedDebitsTotal
	assert.Equal(t, expectedOverallTotal, summary.Overall.TotalAmount)
}
