package main

import (
	"fmt"
	"log"
	"transaction-summary/src/models/libs"
)

func main() {
	csvReader := libs.NewCsvReader()
	transactions, err := csvReader.LoadTransactions("assets/transactions.csv")
	if err != nil {
		log.Fatal(err)
	}

	for _, transaction := range transactions {
		fmt.Println(transaction.String())
	}
}
