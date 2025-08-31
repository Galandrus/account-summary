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
	service services.TransactionsServiceInterface
}

func NewTransactionsHandler(service services.TransactionsServiceInterface) TransactionsHandler {
	return &transactionsHandler{service: service}
}

func (h *transactionsHandler) GetTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountEmail := req.URL.Query().Get("accountEmail")
	if accountEmail == "" {
		http.Error(res, "accountEmail is required", http.StatusBadRequest)
		return
	}

	transactions, err := h.service.GetTransactionsByEmail(ctx, accountEmail)
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
		Path         string `json:"path"`
		AccountEmail string `json:"accountEmail"`
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
	if data.AccountEmail == "" {
		http.Error(res, "accountEmail is required", http.StatusBadRequest)
		return
	}

	log.Println("Loading transactions from path:", data.Path)

	err = h.service.LoadTransactions(ctx, data.Path, data.AccountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(fmt.Sprintf("Transactions loaded successfully for account %s", data.AccountEmail))
}

func (h *transactionsHandler) GetSummary(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountEmail := req.URL.Query().Get("accountEmail")
	if accountEmail == "" {
		http.Error(res, "accountEmail is required", http.StatusBadRequest)
		return
	}

	summary, err := h.service.GetSummaryByEmail(ctx, accountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(summary.String()))
}

func (h *transactionsHandler) SendEmail(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	accountEmail := req.URL.Query().Get("accountEmail")
	if accountEmail == "" {
		http.Error(res, "accountEmail is required", http.StatusBadRequest)
		return
	}

	// El servicio ahora devuelve el HTML ya procesado
	err := h.service.SendEmail(ctx, accountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Email sent successfully"))
}
