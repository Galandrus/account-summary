package main

import (
	"account-summary/src/config"
	"account-summary/src/connections"
	"account-summary/src/handlers"
	"account-summary/src/pkg/csv"
	"account-summary/src/pkg/email"
	"account-summary/src/pkg/files"
	"account-summary/src/pkg/utils"
	"account-summary/src/repository"
	"account-summary/src/server"
	"account-summary/src/services"
	"context"
)

func main() {
	cfg := config.Load()

	mongoClient := connections.NewMongoConn(cfg.MongoURI)
	defer mongoClient.Disconnect(context.Background())

	transactionRepository := repository.NewTransactionRepository(mongoClient)
	accountRepository := repository.NewAccountRepository(mongoClient)

	loader := files.NewLocalLoader()
	summaryProcessor := utils.NewSummaryProcessor()
	csvReader := csv.NewCsvReader(loader)
	emailSender := email.NewEmailSender(cfg)
	idGenerator := utils.NewIdGenerator()

	accountsService := services.NewAccountsService(accountRepository, idGenerator, emailSender)
	transactionsService := services.NewTransactionsService(transactionRepository, accountsService, csvReader, summaryProcessor)

	handler := handlers.NewMainApiHandler(transactionsService, accountsService)

	server := server.NewServer(cfg, handler)
	server.Start()
}
