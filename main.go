package main

import (
	"account-summary/src/libs"
	"fmt"
	"log"
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

	summaryProcessor := libs.NewSummaryProcessor()
	summary, err := summaryProcessor.ProcessSummary(transactions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(summary.String())
}
