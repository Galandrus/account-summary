package handlers

import "net/http"

type MainApiHandlerInterface interface {
	GetTransactions(res http.ResponseWriter, req *http.Request)
	LoadTransactions(res http.ResponseWriter, req *http.Request)
	GetSummary(res http.ResponseWriter, req *http.Request)
	SendSummaryEmail(res http.ResponseWriter, req *http.Request)
}
