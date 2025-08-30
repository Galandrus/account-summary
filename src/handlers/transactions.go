package handlers

import (
	"account-summary/src/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TransactionsHandler interface {
	GetTransactions(res http.ResponseWriter, req *http.Request)
	LoadTransactions(res http.ResponseWriter, req *http.Request)
	GetSummary(res http.ResponseWriter, req *http.Request)
	SendEmail(res http.ResponseWriter, req *http.Request)
}

type transactionsHandler struct {
	service services.TransactionsService
}

func NewTransactionsHandler(service services.TransactionsService) TransactionsHandler {
	return &transactionsHandler{service: service}
}

func (h *transactionsHandler) GetTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountId := req.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(res, "accountId is required", http.StatusBadRequest)
		return
	}

	transactions, err := h.service.GetTransactions(ctx, accountId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transactions)
}

func (h *transactionsHandler) LoadTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body := req.Body
	defer body.Close()

	var data struct {
		Path      string `json:"path"`
		AccountId string `json:"accountId"`
	}
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Path == "" {
		http.Error(res, "path is required", http.StatusBadRequest)
		return
	}
	if data.AccountId == "" {
		http.Error(res, "accountId is required", http.StatusBadRequest)
		return
	}

	log.Println("Loading transactions from path:", data.Path)

	err = h.service.LoadTransactions(ctx, data.Path, data.AccountId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(fmt.Sprintf("Transactions loaded successfully for account %s", data.AccountId))
}

func (h *transactionsHandler) GetSummary(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountId := req.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(res, "accountId is required", http.StatusBadRequest)
		return
	}

	summary, err := h.service.GetSummary(ctx, accountId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(summary.String()))
}

func (h *transactionsHandler) SendEmail(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountId := req.URL.Query().Get("accountId")
	if accountId == "" {
		http.Error(res, "accountId is required", http.StatusBadRequest)
		return
	}

	// El servicio ahora devuelve el HTML ya procesado
	htmlContent, err := h.service.SendEmail(ctx, accountId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Solo establecer headers y devolver el contenido
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write(htmlContent)
}
