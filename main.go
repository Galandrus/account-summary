package main

import (
	"account-summary/src/connections"
	"account-summary/src/handlers"
	"account-summary/src/libs"
	"account-summary/src/repository"
	"account-summary/src/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// csvReader := libs.NewCsvReader()
	// transactions, err := csvReader.LoadTransactions("assets/transactions.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, transaction := range transactions {
	// 	fmt.Println(transaction.String())
	// }

	mongoClient := connections.NewMongoConn("")

	repository := repository.NewTransactionRepository(mongoClient)
	// errCreate := repository.CreateTransactions(context.Background(), transactions)
	// if errCreate != nil {
	// 	log.Fatal(errCreate)
	// }

	summaryProcessor := libs.NewSummaryProcessor()
	csvReader := libs.NewCsvReader()

	service := services.NewTransactionsService(repository, csvReader, summaryProcessor)
	handler := handlers.NewTransactionsHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/load-transactions", handler.LoadTransactions)
	mux.HandleFunc("/transactions", handler.GetTransactions)
	mux.HandleFunc("/summary", handler.GetSummary)

	fmt.Printf("server running on http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))

	// summaryProcessor := libs.NewSummaryProcessor()
	// summary, err := summaryProcessor.ProcessSummary(transactions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(summary.String())
}
