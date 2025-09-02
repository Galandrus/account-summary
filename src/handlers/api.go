package handlers

import (
	"account-summary/src/interfaces/handlers"
	"account-summary/src/interfaces/services"
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

	sendJSONResponse(res, transactions, http.StatusOK)
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

	sendJSONResponse(res, fmt.Sprintf("Transactions loaded successfully for account %s", data.AccountEmail), http.StatusOK)
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

	sendTextResponse(res, account.Summary.String(), http.StatusOK)
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

	sendTextResponse(res, "Email sent successfully", http.StatusOK)
}
