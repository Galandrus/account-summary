package handlers

import (
	"account-summary/src/internal/handlers"
	"account-summary/src/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type mainApiHandler struct {
	transactionsService services.TransactionsServiceInterface
	accountsService     services.AccountsServiceInterface
}

func NewMainApiHandler(
	transactionsService services.TransactionsServiceInterface,
	accountsService services.AccountsServiceInterface,
) handlers.MainApiHandlerInterface {
	return &mainApiHandler{
		transactionsService: transactionsService,
		accountsService:     accountsService,
	}
}

func (h *mainApiHandler) GetTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountEmail, errGetAccountEmail := getAccountEmailQueryParam(req.URL.Query())
	if errGetAccountEmail != nil {
		http.Error(res, errGetAccountEmail.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := h.transactionsService.GetTransactionsByEmail(ctx, accountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transactions)
}

func (h *mainApiHandler) LoadTransactions(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := getLoadTransactionsBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.transactionsService.LoadTransactions(ctx, data.Path, data.AccountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(fmt.Sprintf("Transactions loaded successfully for account %s", data.AccountEmail))
}

func (h *mainApiHandler) GetSummary(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountEmail, errGetAccountEmail := getAccountEmailQueryParam(req.URL.Query())
	if errGetAccountEmail != nil {
		http.Error(res, errGetAccountEmail.Error(), http.StatusBadRequest)
		return
	}

	account, err := h.accountsService.GetAccountByEmail(ctx, accountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(account.Summary.String()))
}

func (h *mainApiHandler) SendSummaryEmail(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := getSendEmailBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.accountsService.SendSummaryEmail(ctx, data.AccountEmail)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Email sent successfully"))
}
