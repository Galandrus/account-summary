package libs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"transaction-summary/src/models"
)

type CsvReader interface {
	LoadTransactions(path string) ([]models.Transaction, error)
}

type csvReader struct{}

func NewCsvReader() CsvReader {
	return &csvReader{}
}

func (c *csvReader) LoadTransactions(path string) ([]models.Transaction, error) {
	fileCsv, err := c.loadFile(path)
	if err != nil {
		return nil, err
	}
	defer fileCsv.Close()

	return c.readCsvFile(fileCsv)
}

func (c *csvReader) loadFile(path string) (*os.File, error) {
	fileCsv, err := os.Open(path)
	if err != nil {
		log.Default().Printf("Error to open csv file: %v\n", err)
		return nil, fmt.Errorf("error to open csv file: %v", err)
	}

	return fileCsv, nil
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
		if len(t) != 4 {
			log.Default().Printf("error to read transaction: %v\n", t)
			continue
		}

		transaction := models.Transaction{
			ID:          t[0],
			Date:        t[1],
			Transaction: t[2],
			Name:        t[3],
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
