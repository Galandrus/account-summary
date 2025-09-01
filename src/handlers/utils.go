package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type loadTransactionsBody struct {
	Path         string `json:"path"`
	AccountEmail string `json:"accountEmail"`
}

type sendEmailBody struct {
	AccountEmail string `json:"accountEmail"`
}

func getAccountEmailQueryParam(query url.Values) (string, error) {
	accountEmail := query.Get("accountEmail")
	if accountEmail == "" {
		return "", errors.New("accountEmail is required")
	}

	return accountEmail, nil
}

func getSendEmailBody(req *http.Request) (sendEmailBody, error) {
	data, err := getBody[sendEmailBody](req)
	if err != nil {
		return sendEmailBody{}, err
	}

	if data.AccountEmail == "" {
		return sendEmailBody{}, errors.New("accountEmail is required")
	}

	return *data, nil

}

func getLoadTransactionsBody(req *http.Request) (loadTransactionsBody, error) {
	data, err := getBody[loadTransactionsBody](req)
	if err != nil {
		return loadTransactionsBody{}, err
	}

	if data.Path == "" {
		return loadTransactionsBody{}, errors.New("path is required")
	}

	if data.AccountEmail == "" {
		return loadTransactionsBody{}, errors.New("accountEmail is required")
	}

	return *data, nil
}

func getBody[T any](req *http.Request) (*T, error) {
	body := req.Body
	defer body.Close()

	var data T
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func sendJSONResponse(res http.ResponseWriter, data interface{}, statusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
}

func sendTextResponse(res http.ResponseWriter, data string, statusCode int) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(statusCode)
	res.Write([]byte(data))
}
