package csv

import (
	"account-summary/src/internal/libs"
	"account-summary/src/models"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type csvReader struct {
	loader libs.FileLoaderInterface
}

func NewCsvReader(loader libs.FileLoaderInterface) libs.CsvReaderInterface {
	return &csvReader{loader: loader}
}

func (c *csvReader) LoadTransactions(path string) ([]models.Transaction, error) {
	fileCsv, err := c.loader.LoadFile(path)
	if err != nil {
		return nil, err
	}
	defer fileCsv.Close()

	return c.readCsvFile(fileCsv)
}

func (c *csvReader) readCsvFile(fileCsv *os.File) ([]models.Transaction, error) {
	readCsv := csv.NewReader(fileCsv)
	readCsv.LazyQuotes = true
	readCsv.Comma = ','

	csvRead, err := readCsv.ReadAll()
	if err != nil {
		msg := fmt.Sprintf("error to read csv file: %v", err)
		log.Default().Println(msg)
		return nil, errors.New(msg)
	}

	transactions, err := c.mapCsvToTransactions(csvRead)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (c *csvReader) mapCsvToTransactions(csvRead [][]string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	csvRead = csvRead[1:]
	for _, t := range csvRead {
		transaction, err := getTransaction(t)
		if err != nil {
			log.Default().Printf("error to read transaction: %v\n", err)
			continue
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func getTransaction(t []string) (models.Transaction, error) {
	if len(t) != 4 {
		return models.Transaction{}, fmt.Errorf("transaction is not valid: %v", t)
	}

	amount, err := strconv.ParseFloat(t[2], 64)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error to parse amount: %v", err)
	}

	date, err := time.Parse("2006-01-02", t[1])
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error to parse date: %v", err)
	}

	transaction := models.Transaction{
		ID:     t[0],
		Date:   date,
		Amount: amount,
		Name:   t[3],
	}

	return transaction, nil
}
