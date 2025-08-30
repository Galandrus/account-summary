package handlers

import (
	"account-summary/src/services"
	"encoding/json"
	"log"
	"net/http"
)

type TransactionsHandler interface {
	GetTransactions(res http.ResponseWriter, req *http.Request)
	LoadTransactions(res http.ResponseWriter, req *http.Request)
	GetSummary(res http.ResponseWriter, req *http.Request)
}

type transactionsHandler struct {
	service services.TransactionsService
}

func NewTransactionsHandler(service services.TransactionsService) TransactionsHandler {
	return &transactionsHandler{service: service}
}

func (h *transactionsHandler) GetTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	transactions, err := h.service.GetTransactions(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(transactions)
}

func (h *transactionsHandler) LoadTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	path := req.URL.Query().Get("path")
	if path == "" {
		http.Error(res, "path is required", http.StatusBadRequest)
		return
	}

	log.Println("Loading transactions from path:", path)

	err := h.service.LoadTransactions(ctx, path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode("OK")
}

func (h *transactionsHandler) GetSummary(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	summary, err := h.service.GetSummary(ctx)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(summary.String()))
}
